package dedash

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const emDash = "\u2014"

type Options struct {
	Root   string
	DryRun bool
}

type FileChange struct {
	Path         string
	Replacements int
}

type Report struct {
	Scanned  int
	Modified int
	Skipped  int
	Changes  []FileChange
}

func Run(ctx context.Context, opts Options) (*Report, error) {
	_ = ctx

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

	report := &Report{}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if path != root && shouldSkipDir(d.Name()) {
				slog.Debug("skipping directory", "path", path)
				return filepath.SkipDir
			}
			return nil
		}

		if !d.Type().IsRegular() {
			slog.Debug("skipping non-regular file", "path", path)
			report.Skipped++
			return nil
		}

		if shouldSkipExtension(d.Name()) {
			slog.Debug("skipping extension", "path", path)
			report.Skipped++
			return nil
		}

		report.Scanned++

		change, err := processFile(path, opts.DryRun)
		if err != nil {
			slog.Debug("skipping file", "path", path, "error", err)
			report.Skipped++
			return nil
		}

		if change.Replacements == 0 {
			return nil
		}

		report.Modified++
		relPath, relErr := filepath.Rel(root, path)
		if relErr != nil {
			relPath = path
		}
		change.Path = relPath
		report.Changes = append(report.Changes, *change)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}

	return report, nil
}

func processFile(path string, dryRun bool) (*FileChange, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if isBinaryContent(data) {
		return nil, fmt.Errorf("binary content")
	}

	content := string(data)
	if !strings.Contains(content, emDash) {
		return &FileChange{Path: path, Replacements: 0}, nil
	}

	replaced := strings.ReplaceAll(content, emDash, "-")
	count := strings.Count(content, emDash)

	if !dryRun {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(replaced), info.Mode()); err != nil {
			return nil, err
		}
	}

	return &FileChange{Path: path, Replacements: count}, nil
}

func TotalReplacements(report *Report) int {
	total := 0
	for _, change := range report.Changes {
		total += change.Replacements
	}
	return total
}
