package kube

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFakeKubectl(t *testing.T, dir, currentCtx string) string {
	t.Helper()

	path := filepath.Join(dir, "kubectl")
	script := fmt.Sprintf(`#!/bin/sh
case "$1" in
config)
  case "$2" in
  current-context) echo %q; exit 0 ;;
  use-context) exit 0 ;;
  esac ;;
port-forward) exit 0 ;;
esac
exit 1
`, currentCtx)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCLIBinDefault(t *testing.T) {
	c := CLI{}
	if got := c.bin(); got != "kubectl" {
		t.Fatalf("bin() = %q, want kubectl", got)
	}
}

func TestCLIBinCustom(t *testing.T) {
	c := CLI{Bin: "/custom/kubectl"}
	if got := c.bin(); got != "/custom/kubectl" {
		t.Fatalf("bin() = %q", got)
	}
}

func TestCurrentContext(t *testing.T) {
	dir := t.TempDir()
	bin := writeFakeKubectl(t, dir, "my-ctx")

	ctx, err := CLI{Bin: bin}.CurrentContext(context.Background())
	if err != nil {
		t.Fatalf("CurrentContext() error = %v", err)
	}
	if ctx != "my-ctx" {
		t.Fatalf("context = %q, want my-ctx", ctx)
	}
}

func TestCurrentContextError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "kubectl")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := CLI{Bin: bin}.CurrentContext(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "kubectl current-context") {
		t.Fatalf("error = %v", err)
	}
}

func TestUseContext(t *testing.T) {
	dir := t.TempDir()
	bin := writeFakeKubectl(t, dir, "ctx-a")

	cli := CLI{Bin: bin}
	if err := cli.UseContext(context.Background(), "ctx-b"); err != nil {
		t.Fatalf("UseContext() error = %v", err)
	}
}

func TestUseContextError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "kubectl")
	script := "#!/bin/sh\nif [ \"$2\" = use-context ]; then echo fail >&2; exit 1; fi\nexit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	err := CLI{Bin: bin}.UseContext(context.Background(), "bad")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "kubectl use-context bad") {
		t.Fatalf("error = %v", err)
	}
}

func TestPortForward(t *testing.T) {
	dir := t.TempDir()
	bin := writeFakeKubectl(t, dir, "ctx")

	cli := CLI{Bin: bin}
	if err := cli.PortForward(context.Background(), []string{"-n", "ns", "svc/x", "8080:80"}); err != nil {
		t.Fatalf("PortForward() error = %v", err)
	}
}

func TestPortForwardError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "kubectl")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	err := CLI{Bin: bin}.PortForward(context.Background(), []string{"8080:80"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "kubectl port-forward") {
		t.Fatalf("error = %v", err)
	}
}
