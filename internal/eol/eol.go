package eol

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Robert27/eggl-cli/internal/git"
	"github.com/Robert27/eggl-cli/internal/ui"
)

// MaxFileSize is the largest regular file eol will read (50 MiB).
const MaxFileSize = 50 << 20

type GitDiffRunner interface {
	InsideWorkTree(ctx context.Context) (bool, error)
	RepoRoot(ctx context.Context) (string, error)
	ChangedFilePaths(ctx context.Context) ([]string, error)
	ChangedFilePathsSince(ctx context.Context, base string) ([]string, error)
}

type Options struct {
	Root          string
	DryRun        bool
	Yes           bool
	Extensions    []string
	IncludeHidden bool
	GitDiff       bool
	GitDiffBase   string
	Git           GitDiffRunner
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
		"git_diff", opts.GitDiff,
		"git_diff_base", opts.GitDiffBase,
	)

	if opts.GitDiff || opts.GitDiffBase != "" {
		if err := scanGitDiff(ctx, opts, root, extensions, report, &pending); err != nil {
			return nil, err
		}
	} else {
		err = walkDirectory(ctx, opts, root, extensions, report, &pending)
		if err != nil {
			return nil, err
		}
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

		prompt := fmt.Sprintf("This will fix %d line endings in %d files. Continue? [y/N]: ", total, report.Modified)
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

func normalizeAbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}

func gitRunner(opts Options) GitDiffRunner {
	if opts.Git != nil {
		return opts.Git
	}
	return git.CLI{}
}

func scanGitDiff(ctx context.Context, opts Options, root string, extensions []string, report *Report, pending *[]pendingChange) error {
	runner := gitRunner(opts)

	inside, err := runner.InsideWorkTree(ctx)
	if err != nil || !inside {
		return fmt.Errorf("not inside a git work tree")
	}

	repoRoot, err := runner.RepoRoot(ctx)
	if err != nil {
		return err
	}
	repoRoot, err = normalizeAbsPath(repoRoot)
	if err != nil {
		return fmt.Errorf("resolve repo root: %w", err)
	}
	scanRoot, err := normalizeAbsPath(root)
	if err != nil {
		return fmt.Errorf("resolve scan root: %w", err)
	}

	var changed []string
	if opts.GitDiffBase != "" {
		changed, err = runner.ChangedFilePathsSince(ctx, opts.GitDiffBase)
	} else {
		changed, err = runner.ChangedFilePaths(ctx)
	}
	if err != nil {
		return err
	}

	slog.Debug("git diff files", "count", len(changed), "base", opts.GitDiffBase)

	for _, rel := range changed {
		if err := ctx.Err(); err != nil {
			return err
		}

		abs := filepath.Join(repoRoot, rel)
		relToRoot, err := filepath.Rel(scanRoot, abs)
		if err != nil || strings.HasPrefix(relToRoot, "..") {
			slog.Debug("skipping file", "path", rel, "reason", "outside scan root")
			continue
		}

		info, err := os.Lstat(abs)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Debug("skipping file", "path", rel, "reason", "deleted file")
				continue
			}
			slog.Debug("skipping file", "path", rel, "reason", err.Error())
			report.Skipped++
			continue
		}
		if !info.Mode().IsRegular() {
			slog.Debug("skipping file", "path", rel, "reason", "not a regular file")
			report.Skipped++
			continue
		}
		if skipReason := skipPath(info.Name(), extensions, opts.IncludeHidden); skipReason != "" {
			slog.Debug("skipping file", "path", rel, "reason", skipReason)
			report.Skipped++
			continue
		}
		if info.Size() > MaxFileSize {
			slog.Debug("skipping file", "path", rel, "reason", "file too large", "size", info.Size())
			report.Skipped++
			continue
		}

		if err := considerFile(ctx, opts, root, abs, relToRoot, extensions, report, pending); err != nil {
			return err
		}
	}

	return nil
}

func walkDirectory(ctx context.Context, opts Options, root string, extensions []string, report *Report, pending *[]pendingChange) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
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

		if skipReason := skipPath(d.Name(), extensions, opts.IncludeHidden); skipReason != "" {
			slog.Debug("skipping file", "path", path, "reason", skipReason)
			report.Skipped++
			return nil
		}

		if fi, err := d.Info(); err == nil && fi.Size() > MaxFileSize {
			slog.Debug("skipping file", "path", path, "reason", "file too large", "size", fi.Size())
			report.Skipped++
			return nil
		}

		relPath, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relPath = path
		}

		return considerFile(ctx, opts, root, path, relPath, extensions, report, pending)
	})
	if err != nil {
		return fmt.Errorf("walk directory: %w", err)
	}
	return nil
}

func skipPath(name string, extensions []string, includeHidden bool) string {
	if !includeHidden && isHiddenName(name) {
		return "hidden path"
	}
	if shouldSkipFile(name) {
		return "ignored file name"
	}
	if !matchesExtension(name, extensions) {
		return "extension not included"
	}
	return ""
}

func considerFile(ctx context.Context, opts Options, root, absPath, relPath string, extensions []string, report *Report, pending *[]pendingChange) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	report.Scanned++

	change, content, mode, err := analyzeFile(absPath)
	if err != nil {
		slog.Debug("skipping file", "path", absPath, "reason", err.Error())
		report.Skipped++
		return nil
	}

	if change.Replacements == 0 {
		return nil
	}

	change.Path = relPath

	report.Modified++
	report.Changes = append(report.Changes, change)
	*pending = append(*pending, pendingChange{
		absPath: absPath,
		content: content,
		mode:    mode,
		change:  change,
	})
	slog.Debug("found line endings to fix",
		"path", relPath,
		"replacements", change.Replacements,
		"dry_run", opts.DryRun,
	)
	return nil
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

	normalized, count := normalizeEOL(data)
	if count == 0 {
		return FileChange{Path: path, Replacements: 0}, nil, info.Mode(), nil
	}

	return FileChange{Path: path, Replacements: count}, normalized, info.Mode(), nil
}

func normalizeEOL(data []byte) ([]byte, int) {
	if !bytes.Contains(data, []byte{'\r'}) {
		return nil, 0
	}

	s := string(data)
	crlf := strings.Count(s, "\r\n")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	cr := strings.Count(s, "\r")
	if cr > 0 {
		s = strings.ReplaceAll(s, "\r", "\n")
	}

	total := crlf + cr
	if total == 0 {
		return nil, 0
	}
	return []byte(s), total
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
