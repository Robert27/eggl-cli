package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func runHelp(t *testing.T, args ...string) string {
	t.Helper()

	rootCmd.SetArgs(args)
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	resetCommandFlags(rootCmd)
	return out.String()
}

func resetCommandFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetCommandFlags(sub)
	}
}

func TestRootHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "--help")

	for _, want := range []string{
		"Available Commands:",
		"dedash",
		"doctor",
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
		"dedash",
		"Flags:",
		"--path",
		"--dry-run",
		"Global Flags:",
		"--verbose",
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
