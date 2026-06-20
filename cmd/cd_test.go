package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCDHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "cd")

	for _, want := range []string{
		"Print a configured directory path",
		"list         List configured directories",
		"--config",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in cd help, got %q", want, output)
		}
	}
}

func TestCDResolve(t *testing.T) {
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

	stdout, _, err := runCmd(t, "cd", "--config", path, "homelab")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if strings.TrimSpace(stdout) != target {
		t.Fatalf("stdout = %q, want %q", stdout, target)
	}
}

func TestCDList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
directories:
  homelab: ~/projects/homelab
  work: /tmp/work
`), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "cd", "--config", path, "list")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "homelab") || !strings.Contains(stdout, "work") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCDListEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
`), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "cd", "--config", path, "list")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "no directories") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCDUnknown(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
directories:
  homelab: ~/projects/homelab
`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := runCmd(t, "cd", "--config", path, "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown directory") {
		t.Fatalf("error = %v", err)
	}
}

func TestCDCompleteMissingConfig(t *testing.T) {
	cdConfigPath = filepath.Join(t.TempDir(), "missing.yaml")
	t.Cleanup(func() { cdConfigPath = "" })

	names, directive := cdComplete(cdCmd, []string{}, "")
	if len(names) != 0 {
		t.Fatalf("names = %v, want none", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v", directive)
	}
}

func TestCDComplete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
directories:
  homelab: ~/projects/homelab
  work: /tmp/work
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cdConfigPath = path
	t.Cleanup(func() { cdConfigPath = "" })

	names, directive := cdComplete(cdCmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v", directive)
	}
	if len(names) != 2 {
		t.Fatalf("names = %v, want 2 entries", names)
	}

	names, directive = cdComplete(cdCmd, []string{"homelab"}, "")
	if len(names) != 0 {
		t.Fatalf("expected no completions with arg, got %v", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v", directive)
	}
}
