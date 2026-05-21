package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"
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
