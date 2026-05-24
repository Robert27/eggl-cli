package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionShort(t *testing.T) {
	stdout, _, err := runCmd(t, "version", "--short")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := strings.TrimSpace(stdout)
	if got == "" {
		t.Fatal("expected non-empty version output")
	}
}

func TestVersion(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stdout, _, err := runCmd(t, "version")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	for _, want := range []string{
		"eggl version",
		"commit:",
		"built:",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("expected %q in version output, got %q", want, stdout)
		}
	}
}

func writeFakeDoctorBin(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func doctorPATH(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"kubectl", "git", "tailscale", "netbird"} {
		writeFakeDoctorBin(t, dir, name)
	}
	return dir
}

func TestDoctor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("PATH", doctorPATH(t))
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout, _, err := runCmd(t, "doctor")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout, "[ok] go:") {
		t.Fatalf("expected go check in output, got %q", stdout)
	}
	for _, name := range []string{"kubectl", "git", "tailscale", "netbird"} {
		if !strings.Contains(stdout, "[ok] "+name+":") {
			t.Fatalf("expected %s check in output, got %q", name, stdout)
		}
	}
	if !strings.Contains(stdout, "All checks passed") {
		t.Fatalf("expected success summary in output, got %q", stdout)
	}
}

func TestDoctorFailure(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("PATH", doctorPATH(t))
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	if err := os.WriteFile(path, []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "doctor", "--check-path", path)
	if err == nil {
		t.Fatal("expected error for non-directory check path")
	}
	if !strings.Contains(stdout, "[fail] home:") {
		t.Fatalf("expected failed home check in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "1 check(s) failed") {
		t.Fatalf("expected failure summary in output, got %q", stdout)
	}
}

func TestVerboseLogging(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("plain text"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, stderr, err := runCmd(t, "--verbose", "dedash", "--path", root, "--dry-run")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stderr, "running dedash") {
		t.Fatalf("expected debug log in stderr, got %q", stderr)
	}
}
