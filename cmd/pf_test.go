package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
