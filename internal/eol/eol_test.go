package eol

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func defaultOpts(root string) Options {
	return Options{Root: root, Yes: true}
}

func writeTestFile(t *testing.T, root, name string, content []byte) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestNormalizeEOL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		count int
	}{
		{name: "lf only", input: "hello\nworld\n", want: "hello\nworld\n", count: 0},
		{name: "crlf", input: "hello\r\nworld\r\n", want: "hello\nworld\n", count: 2},
		{name: "cr only", input: "hello\rworld\r", want: "hello\nworld\n", count: 2},
		{name: "mixed", input: "a\r\nb\nc\r", want: "a\nb\nc\n", count: 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, count := normalizeEOL([]byte(tc.input))
			if count != tc.count {
				t.Fatalf("count = %d, want %d", count, tc.count)
			}
			if tc.count == 0 {
				if got != nil {
					t.Fatalf("got = %q, want nil", got)
				}
				return
			}
			if string(got) != tc.want {
				t.Fatalf("got = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRunConvertsCRLF(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "readme.md", []byte("hello\r\nworld\r\n"))

	report, err := Run(context.Background(), defaultOpts(root))
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 1 || TotalReplacements(report) != 2 {
		t.Fatalf("report = %+v", report)
	}

	got, err := os.ReadFile(filepath.Join(root, "readme.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\nworld\n" {
		t.Fatalf("file content = %q", got)
	}
}

func TestRunDryRunLeavesFilesUnchanged(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "readme.md", []byte("hello\r\nworld\r\n"))

	report, err := Run(context.Background(), Options{Root: root, DryRun: true})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 1 {
		t.Fatalf("Modified = %d, want 1", report.Modified)
	}

	got, err := os.ReadFile(filepath.Join(root, "readme.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello\r\nworld\r\n" {
		t.Fatalf("file should be unchanged, got %q", got)
	}
}

func TestRunSkipsBinary(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "data.bin", []byte("hello\x00\r\nworld"))

	report, err := Run(context.Background(), defaultOpts(root))
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 0 || report.Skipped != 1 {
		t.Fatalf("report = %+v", report)
	}
}

type fakeGitDiff struct {
	inside      bool
	root        string
	files       []string
	filesSince  map[string][]string
	insideE     error
	rootE       error
	filesE      error
	filesSinceE error
}

func (f *fakeGitDiff) InsideWorkTree(context.Context) (bool, error) {
	return f.inside, f.insideE
}

func (f *fakeGitDiff) RepoRoot(context.Context) (string, error) {
	return f.root, f.rootE
}

func (f *fakeGitDiff) ChangedFilePaths(context.Context) ([]string, error) {
	return f.files, f.filesE
}

func (f *fakeGitDiff) ChangedFilePathsSince(_ context.Context, base string) ([]string, error) {
	if f.filesSinceE != nil {
		return nil, f.filesSinceE
	}
	return f.filesSince[base], nil
}

func TestRunGitDiffOnlyChangedFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "changed.md", []byte("a\r\n"))
	writeTestFile(t, root, "unchanged.md", []byte("b\r\n"))

	report, err := Run(context.Background(), Options{
		Root:    root,
		Yes:     true,
		GitDiff: true,
		Git: &fakeGitDiff{
			inside: true,
			root:   root,
			files:  []string{"changed.md"},
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 1 {
		t.Fatalf("Modified = %d, want 1", report.Modified)
	}
	if report.Changes[0].Path != "changed.md" {
		t.Fatalf("change path = %q", report.Changes[0].Path)
	}

	got, err := os.ReadFile(filepath.Join(root, "unchanged.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "b\r\n" {
		t.Fatalf("unchanged file should stay CRLF, got %q", got)
	}
}

func TestRunGitDiffBaseOnlyBranchFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "feature.md", []byte("feature\r\n"))
	writeTestFile(t, root, "shared.md", []byte("shared\r\n"))

	report, err := Run(context.Background(), Options{
		Root:        root,
		Yes:         true,
		GitDiffBase: "main",
		Git: &fakeGitDiff{
			inside: true,
			root:   root,
			filesSince: map[string][]string{
				"main": {"feature.md"},
			},
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 1 || report.Changes[0].Path != "feature.md" {
		t.Fatalf("report = %+v", report)
	}
}

func TestRunGitDiffNotInsideWorkTree(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Root:    t.TempDir(),
		GitDiff: true,
		Git: &fakeGitDiff{
			inside: false,
		},
	})
	if err == nil || !strings.Contains(err.Error(), "not inside a git work tree") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunNotTerminalRequiresYes(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "readme.md", []byte("hello\r\n"))

	_, err := Run(context.Background(), Options{
		Root:  root,
		Input: strings.NewReader(""),
	})
	if err == nil {
		t.Fatal("expected error when stdin is not a terminal")
	}
	if !strings.Contains(err.Error(), "not a terminal") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunGitDiffChangedPathsError(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Root:    t.TempDir(),
		GitDiff: true,
		Git: &fakeGitDiff{
			inside: true,
			filesE: errors.New("diff failed"),
		},
	})
	if err == nil || !strings.Contains(err.Error(), "diff failed") {
		t.Fatalf("error = %v", err)
	}
}
