package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDedashDryRun(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello — world"), 0o644); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"dedash", "--path", root, "--dry-run"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "dry-run:") {
		t.Fatalf("expected dry-run prefix in output, got %q", output)
	}
	if !strings.Contains(output, "would modify 1") {
		t.Fatalf("expected would modify summary in output, got %q", output)
	}
	if !strings.Contains(output, "readme.md (1)") {
		t.Fatalf("expected file change in output, got %q", output)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello — world" {
		t.Fatalf("file should be unchanged in dry-run, got %q", got)
	}
}
