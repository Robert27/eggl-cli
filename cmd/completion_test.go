package cmd

import (
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {
	stdout, _, err := runCmd(t, "completion", "bash")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "eggl") {
		t.Fatalf("expected bash completion to reference eggl, got %q", stdout)
	}
}

func TestCompletionZsh(t *testing.T) {
	stdout, _, err := runCmd(t, "completion", "zsh")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "eggl") {
		t.Fatalf("expected zsh completion to reference eggl, got %q", stdout)
	}
}

func TestCompletionFish(t *testing.T) {
	stdout, _, err := runCmd(t, "completion", "fish")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "eggl") {
		t.Fatalf("expected fish completion to reference eggl, got %q", stdout)
	}
}

func TestCompletionPowerShell(t *testing.T) {
	stdout, _, err := runCmd(t, "completion", "powershell")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "eggl") {
		t.Fatalf("expected powershell completion to reference eggl, got %q", stdout)
	}
}
