package pf

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Robert27/eggl-cli/internal/config"
)

type recordingKube struct {
	args []string
}

func (r *recordingKube) CurrentContext(context.Context) (string, error) {
	return "", nil
}

func (r *recordingKube) UseContext(context.Context, string) error {
	return nil
}

func (r *recordingKube) PortForward(_ context.Context, args []string) error {
	r.args = append([]string(nil), args...)
	return nil
}

func testConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	data := `
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
  grafana:
    namespace: monitoring
    resource: svc/grafana
    ports: ["3000:3000"]
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions("/tmp/config.yaml")
	if opts.ConfigPath != "/tmp/config.yaml" {
		t.Fatalf("ConfigPath = %q", opts.ConfigPath)
	}
	if opts.Kube == nil {
		t.Fatal("expected Kube runner")
	}
}

func TestList(t *testing.T) {
	cfg, err := config.Load(testConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	entries := List(cfg)
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}
	if entries[0].Name != "grafana" || entries[1].Name != "longhorn" {
		t.Fatalf("order = %q, %q", entries[0].Name, entries[1].Name)
	}
	if entries[1].Ports[0] != defaultPorts {
		t.Fatalf("default ports = %v", entries[1].Ports)
	}
}

func TestRunResolvesService(t *testing.T) {
	path := testConfig(t)
	k := &recordingKube{}

	if err := Run(context.Background(), "grafana", Options{ConfigPath: path, Kube: k}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"-n", "monitoring", "svc/grafana", "3000:3000"}
	if strings.Join(k.args, " ") != strings.Join(want, " ") {
		t.Fatalf("args = %v, want %v", k.args, want)
	}
}

func TestRunUnknownService(t *testing.T) {
	path := testConfig(t)
	err := Run(context.Background(), "missing", Options{ConfigPath: path, Kube: &recordingKube{}})
	if err == nil || !strings.Contains(err.Error(), "unknown port-forward") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunDefaultPorts(t *testing.T) {
	path := testConfig(t)
	k := &recordingKube{}

	if err := Run(context.Background(), "longhorn", Options{ConfigPath: path, Kube: k}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if k.args[len(k.args)-1] != defaultPorts {
		t.Fatalf("last arg = %q, want %q", k.args[len(k.args)-1], defaultPorts)
	}
}
