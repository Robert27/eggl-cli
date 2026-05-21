package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionShort(t *testing.T) {
	rootCmd.SetArgs([]string{"version", "--short"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := strings.TrimSpace(out.String())
	if got == "" {
		t.Fatal("expected non-empty version output")
	}
}

func TestDoctor(t *testing.T) {
	rootCmd.SetArgs([]string{"doctor"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "[ok] go:") {
		t.Fatalf("expected go check in output, got %q", out.String())
	}

	if !strings.Contains(out.String(), "All checks passed") {
		t.Fatalf("expected success summary in output, got %q", out.String())
	}
}
