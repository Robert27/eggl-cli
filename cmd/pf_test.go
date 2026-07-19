package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPFHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "pf")

	for _, want := range []string{
		"Port-forward configured Kubernetes services",
		"list         List configured port-forwards",
		"--config",
		"--open",
		"-o",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in pf help, got %q", want, output)
		}
	}
}

func TestPFRun(t *testing.T) {
	binDir := t.TempDir()
	kubectl := filepath.Join(binDir, "kubectl")
	if err := os.WriteFile(kubectl, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
    ports: ["8080:80"]
`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := runCmd(t, "pf", "--config", path, "longhorn")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestPFList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
    ports: ["8080:80"]
`), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "pf", "--config", path, "list")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "longhorn") || !strings.Contains(stdout, "longhorn-system") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestPFListEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
`), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout, _, err := runCmd(t, "pf", "--config", path, "list")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "no port_forwards") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestPFComplete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
    ports: ["8080:80"]
  grafana:
    namespace: monitoring
    resource: svc/grafana
    ports: ["3000:3000"]
`), 0o644); err != nil {
		t.Fatal(err)
	}

	pfConfigPath = path
	t.Cleanup(func() { pfConfigPath = "" })

	names, directive := pfComplete(pfCmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v", directive)
	}
	if len(names) != 2 {
		t.Fatalf("names = %v, want 2 entries", names)
	}
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}
	for _, want := range []string{"longhorn", "grafana"} {
		if !nameSet[want] {
			t.Fatalf("missing completion for %q", want)
		}
	}

	names, directive = pfComplete(pfCmd, []string{"longhorn"}, "")
	if len(names) != 0 {
		t.Fatalf("expected no completions with arg, got %v", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v", directive)
	}
}
