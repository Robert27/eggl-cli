package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvPath(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stdout, _, err := runCmd(t, "env", "path")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "config.yaml") {
		t.Fatalf("expected config path in output, got %q", stdout)
	}
}

func TestEnvInit(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	root := t.TempDir()
	path := filepath.Join(root, "eggl", "config.yaml")

	stdout, _, err := runCmd(t, "env", "--config", path, "init")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "Created") {
		t.Fatalf("expected Created in output, got %q", stdout)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file missing: %v", err)
	}
}
