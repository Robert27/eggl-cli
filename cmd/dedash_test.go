package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDedashDryRun(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "dedash", "--path", root, "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout, "dry-run:") {
		t.Fatalf("expected dry-run prefix in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "would modify 1") {
		t.Fatalf("expected would modify summary in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "readme.md (1)") {
		t.Fatalf("expected file change in output, got %q", stdout)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello \u2014 world" {
		t.Fatalf("file should be unchanged in dry-run, got %q", got)
	}
}

func TestDedashExtFilter(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	for name, content := range map[string]string{
		"readme.md": "hello \u2014 world",
		"notes.txt": "plain text",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	stdout, _, err := runCmd(t, "dedash", "--path", root, "--ext", "md", "--yes")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout, "modified 1") {
		t.Fatalf("expected modified summary in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "readme.md (1)") {
		t.Fatalf("expected readme.md in output, got %q", stdout)
	}
	if strings.Contains(stdout, "notes.txt") {
		t.Fatalf("expected notes.txt to be excluded, got %q", stdout)
	}

	got, err := os.ReadFile(filepath.Join(root, "notes.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "plain text" {
		t.Fatalf("notes.txt should be unchanged, got %q", got)
	}
}

func TestDedashWrites(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "dedash", "--path", root, "--yes")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout, "modified 1") {
		t.Fatalf("expected modified summary in output, got %q", stdout)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello - world" {
		t.Fatalf("file content = %q, want %q", got, "hello - world")
	}
}

func TestDedashRequiresYesWithoutTTY(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, stderr, err := runCmdWithIn(t, strings.NewReader(""), "dedash", "--path", root)
	if err == nil {
		t.Fatal("expected error when stdin is not a terminal and --yes is not set")
	}
	if !strings.Contains(stderr, "not a terminal") && !strings.Contains(err.Error(), "not a terminal") {
		t.Fatalf("error = %v, stderr = %q", err, stderr)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello \u2014 world" {
		t.Fatalf("file should be unchanged, got %q", got)
	}
}

func TestDedashDiffDryRun(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	t.Setenv("NO_COLOR", "1")

	dir := setupDedashGitRepo(t)
	writeDedashFile(t, dir, "readme.md", "original\n")
	runDedashGit(t, dir, "add", "readme.md")
	runDedashGit(t, dir, "commit", "-m", "init")
	writeDedashFile(t, dir, "readme.md", "hello \u2014 world\n")

	stdout, _, err := runCmd(t, "dedash", "--path", dir, "--diff", "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "would modify 1") {
		t.Fatalf("expected would modify summary, got %q", stdout)
	}
	if !strings.Contains(stdout, "readme.md (1)") {
		t.Fatalf("expected readme.md in output, got %q", stdout)
	}

	got, err := os.ReadFile(filepath.Join(dir, "readme.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello \u2014 world\n" {
		t.Fatalf("file should be unchanged in dry-run, got %q", got)
	}
}

func TestDedashDiffBaseDryRun(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	t.Setenv("NO_COLOR", "1")

	dir := setupDedashGitRepo(t)
	writeDedashFile(t, dir, "shared.md", "shared\n")
	runDedashGit(t, dir, "add", "shared.md")
	runDedashGit(t, dir, "commit", "-m", "init")
	runDedashGit(t, dir, "branch", "-M", "main")

	writeDedashFile(t, dir, "feature.md", "feature \u2014 dash\n")
	runDedashGit(t, dir, "checkout", "-b", "feature")
	runDedashGit(t, dir, "add", "feature.md")
	runDedashGit(t, dir, "commit", "-m", "feature")

	stdout, _, err := runCmd(t, "dedash", "--path", dir, "--diff-base", "main", "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "would modify 1") {
		t.Fatalf("expected would modify summary, got %q", stdout)
	}
	if !strings.Contains(stdout, "feature.md (1)") {
		t.Fatalf("expected feature.md in output, got %q", stdout)
	}
	if strings.Contains(stdout, "shared.md") {
		t.Fatalf("expected shared.md to be excluded, got %q", stdout)
	}
}

func TestDedashDiffMutuallyExclusiveWithDiffBase(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	_, _, err := runCmd(t, "dedash", "--diff", "--diff-base", "main", "--dry-run")
	if err == nil {
		t.Fatal("expected error when using --diff with --diff-base")
	}
	if !strings.Contains(err.Error(), "cannot use --diff with --diff-base") {
		t.Fatalf("error = %v", err)
	}
}

func TestDedashDiffOutsideGitRepo(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	writeDedashFile(t, root, "readme.md", "hello \u2014 world\n")

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	_, _, err = runCmd(t, "dedash", "--path", root, "--diff", "--dry-run")
	if err == nil {
		t.Fatal("expected error outside git work tree")
	}
	if !strings.Contains(err.Error(), "not inside a git work tree") {
		t.Fatalf("error = %v", err)
	}
}

func setupDedashGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runDedashGit(t, dir, "init")
	runDedashGit(t, dir, "config", "user.email", "test@example.com")
	runDedashGit(t, dir, "config", "user.name", "Test User")

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

	return dir
}

func runDedashGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
}

func writeDedashFile(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
