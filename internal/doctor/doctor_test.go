package doctor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunUsesCheckPath(t *testing.T) {
	dir := t.TempDir()

	report, err := Run(context.Background(), Options{CheckPath: dir})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if !home.OK {
		t.Fatalf("home check should pass, got %+v", home)
	}
	if home.Status != dir {
		t.Fatalf("home Status = %q, want %q", home.Status, dir)
	}
}

func TestRunDefaultUsesHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	report, err := Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if !home.OK {
		t.Fatalf("home check should pass, got %+v", home)
	}
	if home.Status != dir {
		t.Fatalf("home Status = %q, want %q", home.Status, dir)
	}
}

func TestRunMissingHome(t *testing.T) {
	t.Setenv("HOME", "")

	report, err := Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if home.OK {
		t.Fatal("expected home check to fail when HOME is unset")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunInvalidCheckPath(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")

	report, err := Run(context.Background(), Options{CheckPath: missing})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if home.OK {
		t.Fatal("expected home check to fail for missing path")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunCheckPathNotDirectory(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	if err := os.WriteFile(path, []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{CheckPath: path})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if home.OK {
		t.Fatal("expected home check to fail for non-directory path")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunIncludesRuntimeChecks(t *testing.T) {
	dir := t.TempDir()

	report, err := Run(context.Background(), Options{CheckPath: dir})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for _, name := range []string{"go", "os"} {
		check := findCheck(t, report, name)
		if !check.OK {
			t.Fatalf("%s check should pass, got %+v", name, check)
		}
	}
}

func TestHasFailures(t *testing.T) {
	if HasFailures(&Report{Checks: []Check{{OK: true}}}) {
		t.Fatal("expected no failures")
	}
	if !HasFailures(&Report{Checks: []Check{{OK: false}}}) {
		t.Fatal("expected failures")
	}
}

func findCheck(t *testing.T, report *Report, name string) Check {
	t.Helper()

	for _, check := range report.Checks {
		if check.Name == name {
			return check
		}
	}
	t.Fatalf("check %q not found in report", name)
	return Check{}
}

func writeFakeBin(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestRunToolChecksFound(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"kubectl", "git", "tailscale"} {
		writeFakeBin(t, dir, name)
	}
	t.Setenv("PATH", dir)

	report, err := Run(context.Background(), Options{CheckPath: t.TempDir()})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for _, name := range []string{"kubectl", "git", "tailscale"} {
		check := findCheck(t, report, name)
		if !check.OK {
			t.Fatalf("%s check should pass, got %+v", name, check)
		}
	}
}

func TestRunToolCheckMissing(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	report, err := Run(context.Background(), Options{CheckPath: t.TempDir()})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	kubectl := findCheck(t, report, "kubectl")
	if kubectl.OK {
		t.Fatal("expected kubectl check to fail when not on PATH")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunConfigMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("PATH", t.TempDir())

	report, err := Run(context.Background(), Options{CheckPath: t.TempDir()})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	cfg := findCheck(t, report, "config")
	if !cfg.OK {
		t.Fatalf("config check should pass when missing, got %+v", cfg)
	}
	if cfg.Status != "not found" {
		t.Fatalf("config Status = %q, want not found", cfg.Status)
	}
}

func TestRunConfigInvalid(t *testing.T) {
	dir := t.TempDir()
	configDir := filepath.Join(dir, "eggl")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("profiles: {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("PATH", t.TempDir())

	report, err := Run(context.Background(), Options{CheckPath: t.TempDir()})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	cfg := findCheck(t, report, "config")
	if cfg.OK {
		t.Fatal("expected config check to fail for invalid config")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunConfigValid(t *testing.T) {
	dir := t.TempDir()
	configDir := filepath.Join(dir, "eggl")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	content := `
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("PATH", t.TempDir())

	report, err := Run(context.Background(), Options{CheckPath: t.TempDir()})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	cfg := findCheck(t, report, "config")
	if !cfg.OK {
		t.Fatalf("config check should pass, got %+v", cfg)
	}
	if cfg.Status != configPath {
		t.Fatalf("config Status = %q, want %q", cfg.Status, configPath)
	}
}
