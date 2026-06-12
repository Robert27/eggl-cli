package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func (c CLI) RepoRoot(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, c.bin(), "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("git rev-parse --show-toplevel: empty path")
	}
	return root, nil
}

func (c CLI) ChangedFilePaths(ctx context.Context) ([]string, error) {
	inside, err := c.InsideWorkTree(ctx)
	if err != nil || !inside {
		return nil, fmt.Errorf("not inside a git work tree")
	}

	repoRoot, err := c.RepoRoot(ctx)
	if err != nil {
		return nil, err
	}

	unstaged, err := c.diffNameOnly(ctx)
	if err != nil {
		return nil, err
	}
	staged, err := c.diffNameOnly(ctx, "--cached")
	if err != nil {
		return nil, err
	}

	return existingFiles(repoRoot, uniqueSorted(append(unstaged, staged...))), nil
}

func (c CLI) ChangedFilePathsSince(ctx context.Context, base string) ([]string, error) {
	inside, err := c.InsideWorkTree(ctx)
	if err != nil || !inside {
		return nil, fmt.Errorf("not inside a git work tree")
	}

	if strings.TrimSpace(base) == "" {
		return nil, fmt.Errorf("diff base ref is required")
	}

	repoRoot, err := c.RepoRoot(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := c.resolveCommit(ctx, base); err != nil {
		return nil, err
	}

	paths, err := c.diffNameOnly(ctx, base+"...HEAD")
	if err != nil {
		return nil, err
	}

	return existingFiles(repoRoot, uniqueSorted(paths)), nil
}

func (c CLI) resolveCommit(ctx context.Context, ref string) (string, error) {
	out, err := exec.CommandContext(ctx, c.bin(), "rev-parse", "--verify", ref+"^{commit}").Output()
	if err != nil {
		return "", fmt.Errorf("unknown git ref %q", ref)
	}
	return strings.TrimSpace(string(out)), nil
}

func (c CLI) diffNameOnly(ctx context.Context, extra ...string) ([]string, error) {
	// Only added, copied, modified, renamed, and type-changed paths (not deletions).
	args := append([]string{"diff", "--name-only", "--diff-filter=ACMRT"}, extra...)
	out, err := exec.CommandContext(ctx, c.bin(), args...).Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}

	var paths []string
	for line := range strings.Lines(string(out)) {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths, nil
}

func existingFiles(repoRoot string, paths []string) []string {
	if len(paths) == 0 {
		return nil
	}

	out := make([]string, 0, len(paths))
	for _, rel := range paths {
		if _, err := os.Stat(filepath.Join(repoRoot, rel)); err != nil {
			continue
		}
		out = append(out, rel)
	}
	return out
}

func uniqueSorted(paths []string) []string {
	if len(paths) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(paths))
	unique := make([]string, 0, len(paths))
	for _, path := range paths {
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		unique = append(unique, path)
	}
	sort.Strings(unique)
	return unique
}
