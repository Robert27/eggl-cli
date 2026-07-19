package pf

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/roberteggl/eggl-cli/internal/config"
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

type errKube struct{}

func (e *errKube) CurrentContext(context.Context) (string, error) { return "", nil }
func (e *errKube) UseContext(context.Context, string) error       { return nil }
func (e *errKube) PortForward(context.Context, []string) error {
	return errors.New("port-forward failed")
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
	if opts.ReadyWait == nil {
		t.Fatal("expected ReadyWait")
	}
}

func TestListNil(t *testing.T) {
	if got := List(nil); got != nil {
		t.Fatalf("List(nil) = %v, want nil", got)
	}
}

func TestListEmpty(t *testing.T) {
	cfg := &config.Config{
		Profiles: map[string]config.Profile{
			"a": {KubeContext: "ctx-a", TailscaleAccount: "b3e1"},
		},
	}
	if got := List(cfg); got != nil {
		t.Fatalf("List() = %v, want nil", got)
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
	if entries[0].Namespace != "monitoring" || entries[0].Resource != "svc/grafana" {
		t.Fatalf("grafana entry = %+v", entries[0])
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

func TestRunConfigLoadError(t *testing.T) {
	err := Run(context.Background(), "longhorn", Options{
		ConfigPath: filepath.Join(t.TempDir(), "missing.yaml"),
		Kube:       &recordingKube{},
	})
	if err == nil {
		t.Fatal("expected config load error")
	}
}

func TestRunDefaultConfigPath(t *testing.T) {
	home := filepath.Join(t.TempDir(), "home")
	cfgDir := filepath.Join(home, ".config", "eggl")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cfgDir, "config.yaml")
	data := `
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
port_forwards:
  app:
    namespace: ns
    resource: svc/app
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	k := &recordingKube{}
	if err := Run(context.Background(), "app", Options{Kube: k}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"-n", "ns", "svc/app", defaultPorts}
	if strings.Join(k.args, " ") != strings.Join(want, " ") {
		t.Fatalf("args = %v, want %v", k.args, want)
	}
}

func TestRunNilKube(t *testing.T) {
	binDir := t.TempDir()
	kubectl := filepath.Join(binDir, "kubectl")
	if err := os.WriteFile(kubectl, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)

	path := testConfig(t)
	if err := Run(context.Background(), "longhorn", Options{ConfigPath: path}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunKubePortForwardError(t *testing.T) {
	err := Run(context.Background(), "longhorn", Options{
		ConfigPath: testConfig(t),
		Kube:       &errKube{},
	})
	if err == nil || !strings.Contains(err.Error(), "port-forward failed") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunStderrMessage(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = old
		w.Close()
		r.Close()
	})

	path := testConfig(t)
	done := make(chan struct{})
	var msg string
	go func() {
		defer close(done)
		buf, _ := io.ReadAll(r)
		msg = string(buf)
	}()

	if err := Run(context.Background(), "grafana", Options{ConfigPath: path, Kube: &recordingKube{}}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	w.Close()
	<-done
	if !strings.Contains(msg, "port-forward grafana") || !strings.Contains(msg, "localhost:3000") {
		t.Fatalf("stderr = %q", msg)
	}
}

func TestRunOpen(t *testing.T) {
	path := testConfig(t)
	k := &recordingKube{}
	var opened string
	opts := Options{
		ConfigPath: path,
		Kube:       k,
		Open:       true,
		OpenURL: func(url string) error {
			opened = url
			return nil
		},
		ReadyWait: func(context.Context, string) error { return nil },
	}

	if err := Run(context.Background(), "grafana", opts); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if opened != "http://localhost:3000" {
		t.Fatalf("opened = %q, want http://localhost:3000", opened)
	}
	want := []string{"-n", "monitoring", "svc/grafana", "3000:3000"}
	if strings.Join(k.args, " ") != strings.Join(want, " ") {
		t.Fatalf("args = %v, want %v", k.args, want)
	}
}

func TestRunOpenPortForwardFailsBeforeReady(t *testing.T) {
	path := testConfig(t)
	opts := Options{
		ConfigPath: path,
		Kube:       &errKube{},
		Open:       true,
		OpenURL:    func(string) error { return nil },
		ReadyWait:  waitForLocalPort,
	}

	err := Run(context.Background(), "grafana", opts)
	if err == nil || !strings.Contains(err.Error(), "port-forward failed") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestWaitUntilReadyPortForwardFails(t *testing.T) {
	errCh := make(chan error, 1)
	errCh <- errors.New("port-forward failed")

	err := waitUntilReady(context.Background(), "3000", errCh, func(context.Context, string) error {
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "port-forward failed") {
		t.Fatalf("waitUntilReady() error = %v", err)
	}
}
