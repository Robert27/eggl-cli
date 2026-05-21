package cmd

import (
	"os"
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

func TestDedashWrites(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "dedash", "--path", root)
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
