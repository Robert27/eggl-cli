package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEOLHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "eol", "--help")

	for _, want := range []string{
		"eggl eol",
		"Normalize line endings",
		"--dry-run",
		"--diff",
		"--diff-base",
		"Examples:",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in eol help, got %q", want, output)
		}
	}

	if strings.Contains(output, "Available Commands:") {
		t.Fatalf("eol help should not list subcommands, got %q", output)
	}
}

func TestEOLDryRun(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello\r\nworld\r\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "eol", "--path", root, "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout, "dry-run:") {
		t.Fatalf("expected dry-run prefix in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "would modify 1") {
		t.Fatalf("expected would modify summary in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "readme.md (2)") {
		t.Fatalf("expected file change in output, got %q", stdout)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\r\nworld\r\n" {
		t.Fatalf("file should be unchanged in dry-run, got %q", got)
	}
}

func TestEOLWrites(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello\r\nworld\r\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "eol", "--path", root, "--yes")
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
	if string(got) != "hello\nworld\n" {
		t.Fatalf("file content = %q, want %q", got, "hello\nworld\n")
	}
}

func TestEOLDiffDryRun(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	t.Setenv("NO_COLOR", "1")

	dir := setupEOLGitRepo(t)
	writeEOLFile(t, dir, "readme.md", "original\n")
	runEOLGit(t, dir, "add", "readme.md")
	runEOLGit(t, dir, "commit", "-m", "init")
	writeEOLFile(t, dir, "readme.md", "hello\r\nworld\r\n")

	stdout, _, err := runCmd(t, "eol", "--path", dir, "--diff", "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "would modify 1") {
		t.Fatalf("expected would modify summary, got %q", stdout)
	}
	if !strings.Contains(stdout, "readme.md (2)") {
		t.Fatalf("expected readme.md in output, got %q", stdout)
	}

	got, err := os.ReadFile(filepath.Join(dir, "readme.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\r\nworld\r\n" {
		t.Fatalf("file should be unchanged in dry-run, got %q", got)
	}
}

func TestEOLDiffBaseDryRun(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	t.Setenv("NO_COLOR", "1")

	dir := setupEOLGitRepo(t)
	writeEOLFile(t, dir, "shared.md", "shared\n")
	runEOLGit(t, dir, "add", "shared.md")
	runEOLGit(t, dir, "commit", "-m", "init")
	runEOLGit(t, dir, "branch", "-M", "main")

	writeEOLFile(t, dir, "feature.md", "feature\r\nline\r\n")
	runEOLGit(t, dir, "checkout", "-b", "feature")
	runEOLGit(t, dir, "add", "feature.md")
	runEOLGit(t, dir, "commit", "-m", "feature")

	stdout, _, err := runCmd(t, "eol", "--path", dir, "--diff-base", "main", "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "would modify 1") {
		t.Fatalf("expected would modify summary, got %q", stdout)
	}
	if !strings.Contains(stdout, "feature.md (2)") {
		t.Fatalf("expected feature.md in output, got %q", stdout)
	}
	if strings.Contains(stdout, "shared.md") {
		t.Fatalf("expected shared.md to be excluded, got %q", stdout)
	}
}

func TestEOLDiffMutuallyExclusiveWithDiffBase(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	_, _, err := runCmd(t, "eol", "--diff", "--diff-base", "main", "--dry-run")
	if err == nil {
		t.Fatal("expected error when using --diff with --diff-base")
	}
	if !strings.Contains(err.Error(), "cannot use --diff with --diff-base") {
		t.Fatalf("error = %v", err)
	}
}

func setupEOLGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runEOLGit(t, dir, "init")
	runEOLGit(t, dir, "config", "user.email", "test@example.com")
	runEOLGit(t, dir, "config", "user.name", "Test User")

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

func runEOLGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
}

func writeEOLFile(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
