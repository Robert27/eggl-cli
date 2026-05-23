package dedash

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Robert27/eggl-cli/internal/ui"
)

const emDash = "\u2014"

// MaxFileSize is the largest regular file dedash will read (50 MiB).
const MaxFileSize = 50 << 20

type Options struct {
	Root          string
	DryRun        bool
	Yes           bool
	Extensions    []string
	IncludeHidden bool
	Input         io.Reader
	Output        io.Writer
}

type FileChange struct {
	Path         string
	Replacements int
}

type Report struct {
	Scanned   int
	Modified  int
	Skipped   int
	Changes   []FileChange
	Cancelled bool
}

type pendingChange struct {
	absPath string
	content []byte
	mode    fs.FileMode
	change  FileChange
}

func Run(ctx context.Context, opts Options) (*Report, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	extensions := normalizeExtensions(opts.Extensions)

	report := &Report{}
	var pending []pendingChange

	slog.Debug("scanning directory",
		"root", root,
		"dry_run", opts.DryRun,
		"extensions", extensions,
		"include_hidden", opts.IncludeHidden,
	)

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if path != root {
				if !opts.IncludeHidden && isHiddenName(d.Name()) {
					slog.Debug("skipping directory", "path", path, "reason", "hidden path")
					return filepath.SkipDir
				}
				if shouldSkipDir(d.Name()) {
					slog.Debug("skipping directory", "path", path, "reason", "ignored directory name")
					return filepath.SkipDir
				}
			}
			return nil
		}

		if !d.Type().IsRegular() {
			slog.Debug("skipping file", "path", path, "reason", "not a regular file")
			report.Skipped++
			return nil
		}

		if !opts.IncludeHidden && isHiddenName(d.Name()) {
			slog.Debug("skipping file", "path", path, "reason", "hidden path")
			report.Skipped++
			return nil
		}

		if shouldSkipFile(d.Name()) {
			slog.Debug("skipping file", "path", path, "reason", "ignored file name")
			report.Skipped++
			return nil
		}

		if !matchesExtension(d.Name(), extensions) {
			slog.Debug("skipping file", "path", path, "reason", "extension not included")
			report.Skipped++
			return nil
		}

		if fi, err := d.Info(); err == nil && fi.Size() > MaxFileSize {
			slog.Debug("skipping file", "path", path, "reason", "file too large", "size", fi.Size())
			report.Skipped++
			return nil
		}

		report.Scanned++

		change, content, mode, err := analyzeFile(path)
		if err != nil {
			slog.Debug("skipping file", "path", path, "reason", err.Error())
			report.Skipped++
			return nil
		}

		if change.Replacements == 0 {
			return nil
		}

		relPath, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relPath = path
		}
		change.Path = relPath

		report.Modified++
		report.Changes = append(report.Changes, change)
		pending = append(pending, pendingChange{
			absPath: path,
			content: content,
			mode:    mode,
			change:  change,
		})
		slog.Debug("found em-dashes",
			"path", relPath,
			"replacements", change.Replacements,
			"dry_run", opts.DryRun,
		)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}

	slog.Debug("scan complete",
		"scanned", report.Scanned,
		"modified", report.Modified,
		"skipped", report.Skipped,
	)

	total := TotalReplacements(report)
	if opts.DryRun || total == 0 {
		return report, nil
	}

	if !opts.Yes {
		in := opts.Input
		if in == nil {
			in = os.Stdin
		}
		out := opts.Output
		if out == nil {
			out = os.Stderr
		}
		if !ui.IsInteractiveInput(in) {
			return report, fmt.Errorf("not a terminal; use --yes to confirm writes")
		}

		prompt := fmt.Sprintf("This will replace %d em-dashes in %d files. Continue? [y/N]: ", total, report.Modified)
		ok, err := ui.ConfirmPrompt(out, in, prompt)
		if err != nil {
			return report, err
		}
		if !ok {
			report.Cancelled = true
			return report, nil
		}
	}

	for _, item := range pending {
		if err := writeFile(item.absPath, item.content, item.mode); err != nil {
			return report, err
		}
	}

	return report, nil
}

func analyzeFile(path string) (FileChange, []byte, fs.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileChange{}, nil, 0, err
	}
	if info.Size() > MaxFileSize {
		return FileChange{}, nil, 0, fmt.Errorf("file too large")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return FileChange{}, nil, 0, err
	}

	if isBinaryContent(data) {
		return FileChange{}, nil, 0, fmt.Errorf("binary content")
	}

	content := string(data)
	if !strings.Contains(content, emDash) {
		return FileChange{Path: path, Replacements: 0}, nil, info.Mode(), nil
	}

	replaced := strings.ReplaceAll(content, emDash, "-")
	count := strings.Count(content, emDash)

	return FileChange{Path: path, Replacements: count}, []byte(replaced), info.Mode(), nil
}

func writeFile(path string, content []byte, mode fs.FileMode) error {
	return os.WriteFile(path, content, mode)
}

func TotalReplacements(report *Report) int {
	total := 0
	for _, change := range report.Changes {
		total += change.Replacements
	}
	return total
}
