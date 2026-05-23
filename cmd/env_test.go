package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const envTestAccountsJSON = `[
  {"id":"a7f2","nickname":"user-b","tailnet":"example-beta.internal","account":"user-b@example.com","selected":false},
  {"id":"b3e1","nickname":"user-a","tailnet":"example-alpha.internal","account":"user-a@example.com","selected":true}
]`

func writeEnvTestConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := `profiles:
  alpha:
    kube_context: ctx-a
    tailscale_account: b3e1
  beta:
    kube_context: ctx-b
    tailscale_account: a7f2
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func prependFakeBinaries(t *testing.T, kubeCtx string, tailscaleJSON string) {
	t.Helper()

	dir := t.TempDir()

	kubectl := filepath.Join(dir, "kubectl")
	kubeScript := fmt.Sprintf(`#!/bin/sh
case "$1" in
config)
  case "$2" in
  current-context) echo %q; exit 0 ;;
  use-context) exit 0 ;;
  esac ;;
esac
exit 1
`, kubeCtx)
	if err := os.WriteFile(kubectl, []byte(kubeScript), 0o755); err != nil {
		t.Fatal(err)
	}

	jsonPath := filepath.Join(dir, "accounts.json")
	if err := os.WriteFile(jsonPath, []byte(tailscaleJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	tailscaleBin := filepath.Join(dir, "tailscale")
	tsScript := fmt.Sprintf(`#!/bin/sh
json=%q
if [ "$1" = switch ] && [ "$2" = --list ] && [ "$3" = --json ]; then
  cat "$json"
  exit 0
fi
if [ "$1" = switch ]; then
  exit 0
fi
exit 1
`, jsonPath)
	if err := os.WriteFile(tailscaleBin, []byte(tsScript), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

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

func TestEnvShow(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	path := writeEnvTestConfig(t)
	prependFakeBinaries(t, "ctx-a", envTestAccountsJSON)

	stdout, _, err := runCmd(t, "env", "--config", path, "show")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"profile: alpha", "kube: ctx-a", "tailscale:"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout = %q, want substring %q", stdout, want)
		}
	}
}

func TestEnvToggle(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	path := writeEnvTestConfig(t)
	prependFakeBinaries(t, "ctx-a", envTestAccountsJSON)

	stdout, _, err := runCmd(t, "env", "--config", path, "toggle")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"profile:", "beta", "ctx-b"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout = %q, want substring %q", stdout, want)
		}
	}
}

func TestEnvUse(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	path := writeEnvTestConfig(t)
	prependFakeBinaries(t, "ctx-b", envTestAccountsJSON)

	stdout, _, err := runCmd(t, "env", "--config", path, "use", "alpha")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"profile:", "alpha", "ctx-a"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout = %q, want substring %q", stdout, want)
		}
	}
}
