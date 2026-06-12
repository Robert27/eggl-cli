package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestChangedFilePaths(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	initGitRepo(t, dir)

	writeFile(t, dir, "tracked.md", "original\n")
	runGit(t, dir, "add", "tracked.md")
	runGit(t, dir, "commit", "-m", "init")

	writeFile(t, dir, "tracked.md", "modified \u2014 content\n")
	writeFile(t, dir, "staged.md", "staged \u2014 file\n")
	runGit(t, dir, "add", "staged.md")

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	paths, err := CLI{}.ChangedFilePaths(context.Background())
	if err != nil {
		t.Fatalf("ChangedFilePaths() error = %v", err)
	}

	want := []string{"staged.md", "tracked.md"}
	if len(paths) != len(want) {
		t.Fatalf("ChangedFilePaths() = %v, want %v", paths, want)
	}
	for i, path := range paths {
		if path != want[i] {
			t.Fatalf("ChangedFilePaths()[%d] = %q, want %q", i, path, want[i])
		}
	}
}

func TestChangedFilePathsSinceBranch(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	initGitRepo(t, dir)

	writeFile(t, dir, "shared.md", "shared\n")
	runGit(t, dir, "add", "shared.md")
	runGit(t, dir, "commit", "-m", "init")
	runGit(t, dir, "branch", "-M", "main")

	writeFile(t, dir, "feature.md", "feature \u2014 dash\n")
	runGit(t, dir, "checkout", "-b", "feature")
	runGit(t, dir, "add", "feature.md")
	runGit(t, dir, "commit", "-m", "feature")

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	paths, err := CLI{}.ChangedFilePathsSince(context.Background(), "main")
	if err != nil {
		t.Fatalf("ChangedFilePathsSince() error = %v", err)
	}

	if len(paths) != 1 || paths[0] != "feature.md" {
		t.Fatalf("ChangedFilePathsSince() = %v, want [feature.md]", paths)
	}
}

func TestChangedFilePathsSinceExcludesDeletions(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	initGitRepo(t, dir)

	writeFile(t, dir, "keep.md", "keep\n")
	writeFile(t, dir, "remove.md", "remove\n")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "init")
	runGit(t, dir, "branch", "-M", "main")

	runGit(t, dir, "checkout", "-b", "feature")
	runGit(t, dir, "rm", "remove.md")
	runGit(t, dir, "commit", "-m", "remove file")

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	paths, err := CLI{}.ChangedFilePathsSince(context.Background(), "main")
	if err != nil {
		t.Fatalf("ChangedFilePathsSince() error = %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("ChangedFilePathsSince() = %v, want no existing files", paths)
	}
}

func TestChangedFilePathsSinceUnknownRef(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	initGitRepo(t, dir)
	writeFile(t, dir, "readme.md", "hello\n")
	runGit(t, dir, "add", "readme.md")
	runGit(t, dir, "commit", "-m", "init")

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	_, err = CLI{}.ChangedFilePathsSince(context.Background(), "does-not-exist")
	if err == nil {
		t.Fatal("expected error for unknown ref")
	}
	if !strings.Contains(err.Error(), "unknown git ref") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangedFilePathsExcludesDeletions(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	initGitRepo(t, dir)

	writeFile(t, dir, "keep.md", "keep\n")
	writeFile(t, dir, "remove.md", "remove \u2014 me\n")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "init")

	writeFile(t, dir, "keep.md", "changed \u2014 keep\n")
	runGit(t, dir, "rm", "remove.md")

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	paths, err := CLI{}.ChangedFilePaths(context.Background())
	if err != nil {
		t.Fatalf("ChangedFilePaths() error = %v", err)
	}

	if len(paths) != 1 || paths[0] != "keep.md" {
		t.Fatalf("ChangedFilePaths() = %v, want [keep.md]", paths)
	}
}

func TestExistingFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "present.md", "ok\n")

	got := existingFiles(dir, []string{"present.md", "gone.md"})
	want := []string{"present.md"}
	if len(got) != len(want) {
		t.Fatalf("existingFiles() = %v, want %v", got, want)
	}
}

func TestChangedFilePathsNotInsideWorkTree(t *testing.T) {
	dir := t.TempDir()
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	_, err = CLI{}.ChangedFilePaths(context.Background())
	if err == nil {
		t.Fatal("expected error outside git work tree")
	}
	if !strings.Contains(err.Error(), "not inside a git work tree") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUniqueSorted(t *testing.T) {
	got := uniqueSorted([]string{"b.md", "a.md", "b.md"})
	want := []string{"a.md", "b.md"}
	if len(got) != len(want) {
		t.Fatalf("uniqueSorted() = %v, want %v", got, want)
	}
	for i, path := range got {
		if path != want[i] {
			t.Fatalf("uniqueSorted()[%d] = %q, want %q", i, path, want[i])
		}
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
