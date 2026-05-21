package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestRenderVersionPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderVersion(&buf, VersionInfo{
		Version: "v1.2.3",
		Commit:  "abc123",
		Date:    "2026-01-01",
	})

	got := buf.String()
	for _, want := range []string{
		"eggl version v1.2.3",
		"commit: abc123",
		"built:  2026-01-01",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestRenderDoctorPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderDoctor(&buf, []DoctorCheck{
		{Name: "go", Status: "go1.26", Detail: "Go runtime available", OK: true},
		{Name: "home", Status: "missing", Detail: "HOME not set", OK: false},
	})

	got := buf.String()
	for _, want := range []string{
		"[ok] go: go1.26 (Go runtime available)",
		"[fail] home: missing (HOME not set)",
		"1 check(s) failed",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestRenderDedashPlainDryRun(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderDedash(&buf, DedashSummary{
		Scanned:           3,
		Modified:          1,
		Skipped:           1,
		TotalReplacements: 2,
		DryRun:            true,
		Changes: []DedashChange{
			{Path: "readme.md", Replacements: 2},
		},
	})

	got := buf.String()
	for _, want := range []string{
		"dry-run: scanned 3 files, skipped 1, would modify 1 (2 replacements)",
		"readme.md (2)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestIsInteractiveRespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	if IsInteractive(os.Stdout) {
		t.Fatal("expected non-interactive output when NO_COLOR is set")
	}
}

func TestIsInteractiveNonFileWriter(t *testing.T) {
	var buf bytes.Buffer
	if IsInteractive(&buf) {
		t.Fatal("expected non-interactive for non-*os.File writer")
	}
}

func TestRenderHeaderPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderHeader(&buf, "eggl", "A helper CLI")

	got := buf.String()
	for _, want := range []string{"eggl\n", "A helper CLI\n\n"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestRenderHelpPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.BoolP("verbose", "v", false, "Log operation details")
	flags.String("path", ".", "Directory to scan")

	var buf bytes.Buffer
	RenderHelp(&buf, "A helper CLI", "Longer description\nwith two lines", []HelpCommand{
		{Name: "dedash", Description: "Replace em-dashes"},
		{Name: "doctor", Description: "Environment checks"},
	}, flags)

	got := buf.String()
	for _, want := range []string{
		"A helper CLI\n\n",
		"Description:",
		"Longer description",
		"Available Commands:",
		"dedash",
		"Global Flags:",
		"-v, --verbose",
		"--path",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestRenderCommandHelpPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cmd := &cobra.Command{
		Use:   "dedash",
		Short: "Replace em-dashes with hyphens",
		Long: `Replace em-dashes with hyphens in text files.

Skips binaries and build output.`,
		Example: `  eggl dedash --dry-run
  eggl dedash --path ./docs`,
		ValidArgs: []string{"none"},
	}
	cmd.Flags().String("path", ".", "Directory tree to scan")
	cmd.Flags().Bool("dry-run", false, "Report changes without writing")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Log operation details")

	var buf bytes.Buffer
	RenderCommandHelp(&buf, cmd)

	got := buf.String()
	for _, want := range []string{
		"dedash - Replace em-dashes with hyphens",
		"Usage:",
		"Description:",
		"Skips binaries",
		"Valid Args:",
		"Flags:",
		"--path",
		"--dry-run",
		"-v, --verbose",
		"Examples:",
		"eggl dedash --dry-run",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestCommandDescription(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
		want string
	}{
		{
			name: "long only",
			cmd:  &cobra.Command{Long: "Full long text"},
			want: "Full long text",
		},
		{
			name: "short duplicates first long line",
			cmd:  &cobra.Command{Short: "Same line", Long: "Same line\n\nMore detail"},
			want: "More detail",
		},
		{
			name: "short equals long only",
			cmd:  &cobra.Command{Short: "Only", Long: "Only"},
			want: "",
		},
		{
			name: "empty",
			cmd:  &cobra.Command{},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := CommandDescription(tc.cmd); got != tc.want {
				t.Fatalf("CommandDescription() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRenderDoctorAllPassedPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderDoctor(&buf, []DoctorCheck{
		{Name: "go", Status: "go1.26", Detail: "ok", OK: true},
	})

	if !strings.Contains(buf.String(), "All checks passed") {
		t.Fatalf("expected success summary, got %q", buf.String())
	}
}

func TestRenderDedashNoChangesPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderDedash(&buf, DedashSummary{Scanned: 5, Skipped: 2})

	got := buf.String()
	want := "scanned 5 files, skipped 2, no em-dashes found"
	if !strings.Contains(got, want) {
		t.Fatalf("expected %q in output, got %q", want, got)
	}
}

func TestRenderCommandHelpInheritedFlagsPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	parent := &cobra.Command{Use: "eggl"}
	parent.PersistentFlags().BoolP("verbose", "v", false, "Log operation details")

	child := &cobra.Command{
		Use:   "dedash",
		Short: "Replace em-dashes",
	}
	child.Flags().Bool("dry-run", false, "Preview only")
	parent.AddCommand(child)

	var buf bytes.Buffer
	RenderCommandHelp(&buf, child)

	got := buf.String()
	if !strings.Contains(got, "Global Flags:") {
		t.Fatalf("expected inherited flags section, got %q", got)
	}
	if !strings.Contains(got, "-v, --verbose") {
		t.Fatalf("expected verbose flag in global section, got %q", got)
	}
}

func TestRenderCommandHelpUsageWithoutFlagsInLine(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		DisableFlagsInUseLine: true,
		Short:                 "Generate shell completion",
	}

	var buf bytes.Buffer
	RenderCommandHelp(&buf, cmd)

	if !strings.Contains(buf.String(), "completion [bash|zsh|fish|powershell]") {
		t.Fatalf("expected use line in help, got %q", buf.String())
	}
}

func TestRenderCommandHelpSubtitleFromLong(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cmd := &cobra.Command{
		Use: "doctor",
		Long: `First line only in long.

More detail below.`,
	}

	var buf bytes.Buffer
	RenderCommandHelp(&buf, cmd)

	if !strings.Contains(buf.String(), "doctor - First line only in long.") {
		t.Fatalf("expected subtitle from long, got %q", buf.String())
	}
}

func TestRenderHelpSkipsHiddenFlags(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Bool("help", false, "help")
	flags.MarkHidden("help")
	flags.String("visible", "x", "shown flag")

	var buf bytes.Buffer
	RenderHelp(&buf, "summary", "", nil, flags)

	got := buf.String()
	if strings.Contains(got, "--help") {
		t.Fatalf("hidden help flag should not appear, got %q", got)
	}
	if !strings.Contains(got, "--visible") {
		t.Fatalf("expected visible flag, got %q", got)
	}
}

func TestRenderDedashModifiedPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderDedash(&buf, DedashSummary{
		Scanned:           4,
		Modified:          2,
		TotalReplacements: 3,
		Changes: []DedashChange{
			{Path: "a.md", Replacements: 1},
			{Path: "b.md", Replacements: 2},
		},
	})

	got := buf.String()
	for _, want := range []string{
		"scanned 4 files, modified 2 (3 replacements)",
		"a.md (1)",
		"b.md (2)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}
