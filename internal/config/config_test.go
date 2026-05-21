package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
  b:
    kube_context: ctx-b
    tailscale_account: a7f2
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Profiles) != 2 {
		t.Fatalf("profiles len = %d, want 2", len(cfg.Profiles))
	}
}

func TestValidateRequiresFields(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"bad": {KubeContext: "", TailscaleAccount: "x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestWriteInit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "eggl", "config.yaml")
	if err := WriteInit(path); err != nil {
		t.Fatalf("WriteInit() error = %v", err)
	}
	if err := WriteInit(path); err == nil {
		t.Fatal("expected error when config exists")
	}
}
