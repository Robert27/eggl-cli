package cd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Robert27/eggl-cli/internal/config"
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

func TestResolveAbsolutePath(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "work")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}

	path := writeConfig(t, `
directories:
  work: `+target+`
`)
	got, err := Resolve("work", DefaultOptions(path))
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got != target {
		t.Fatalf("Resolve() = %q, want %q", got, target)
	}
}

func TestResolveInvalidPath(t *testing.T) {
	path := writeConfig(t, `
directories:
  bad: ~other
`)
	_, err := Resolve("bad", DefaultOptions(path))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "directory \"bad\"") {
		t.Fatalf("error = %v", err)
	}
}

func TestResolveMissingConfig(t *testing.T) {
	_, err := Resolve("homelab", DefaultOptions(filepath.Join(t.TempDir(), "missing.yaml")))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListEmpty(t *testing.T) {
	if List(nil) != nil {
		t.Fatal("expected nil for nil config")
	}
	if List(&config.Config{}) != nil {
		t.Fatal("expected nil for empty directories")
	}
}

func TestLoadConfigMissing(t *testing.T) {
	_, path, err := loadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error")
	}
	if path == "" {
		t.Fatal("expected config path in error")
	}
}

func TestLoadConfigDefaultPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	path := filepath.Join(dir, "eggl", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`
directories:
  work: /tmp/work
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, gotPath, err := loadConfig("")
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if gotPath != path {
		t.Fatalf("path = %q, want %q", gotPath, path)
	}
	if cfg.Directories["work"] != "/tmp/work" {
		t.Fatalf("directories = %+v", cfg.Directories)
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
