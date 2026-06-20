package cd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)

	target := filepath.Join(home, "projects", "homelab")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
directories:
  homelab: ~/projects/homelab
`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Resolve("homelab", DefaultOptions(path))
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got != target {
		t.Fatalf("Resolve() = %q, want %q", got, target)
	}
}

func TestResolveUnknown(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
directories:
  homelab: ~/projects/homelab
`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Resolve("missing", DefaultOptions(path))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown directory") {
		t.Fatalf("error = %v", err)
	}
}

func TestList(t *testing.T) {
	path := writeConfig(t, `
directories:
  work: /tmp/work
  homelab: ~/projects/homelab
`)
	cfg, _, err := loadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	entries := List(cfg)
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if entries[0].Name != "homelab" || entries[1].Name != "work" {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions("/tmp/config.yaml")
	if opts.ConfigPath != "/tmp/config.yaml" {
		t.Fatalf("ConfigPath = %q", opts.ConfigPath)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
