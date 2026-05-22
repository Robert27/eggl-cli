package cmd

import (
	"strings"
	"testing"
)

func TestRootHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "--help")

	for _, want := range []string{
		"Available Commands:",
		"dedash",
		"doctor",
		"empty",
		"Global Flags:",
		"--verbose",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in root help, got %q", want, output)
		}
	}
}

func TestDedashHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "dedash", "--help")

	for _, want := range []string{
		"eggl dedash",
		"Usage:",
		"Description:",
		"dedash",
		"Flags:",
		"--path",
		"--dry-run",
		"Examples:",
		"Global Flags:",
		"--verbose",
		"Log operation details to stderr",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in dedash help, got %q", want, output)
		}
	}

	if strings.Contains(output, "Available Commands:") {
		t.Fatalf("dedash help should not list subcommands, got %q", output)
	}
}

func TestDoctorHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "doctor", "--help")

	for _, want := range []string{
		"eggl doctor",
		"--check-path",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in doctor help, got %q", want, output)
		}
	}

	if strings.Contains(output, "Available Commands:") {
		t.Fatalf("doctor help should not list subcommands, got %q", output)
	}
}

func TestEnvHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "env", "--help")

	for _, want := range []string{
		"Available Commands:",
		"toggle",
		"show",
		"use",
		"init",
		"path",
		"--config",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in env help, got %q", want, output)
		}
	}

	if strings.Contains(output, "does not roll back") {
		t.Fatalf("env help should use short description, got %q", output)
	}
}
