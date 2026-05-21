package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBinarySmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping binary smoke test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	root := filepath.Dir(wd)
	bin := filepath.Join(t.TempDir(), "eggl")

	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	out, err := exec.Command(bin, "version", "--short").Output()
	if err != nil {
		t.Fatalf("eggl version --short: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty version output from built binary")
	}
}
