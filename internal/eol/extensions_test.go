package eol

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeExtensions(t *testing.T) {
	got := normalizeExtensions([]string{" md ", "txt,.go", ".MD", "txt"})
	want := []string{".md", ".txt", ".go"}
	if len(got) != len(want) {
		t.Fatalf("normalizeExtensions() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("normalizeExtensions()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestMatchesExtension(t *testing.T) {
	extensions := normalizeExtensions([]string{"md", "txt"})

	if !matchesExtension("readme.md", extensions) {
		t.Fatal("expected .md to match")
	}
	if !matchesExtension("README.MD", extensions) {
		t.Fatal("expected .MD to match case-insensitively")
	}
	if matchesExtension("image.png", extensions) {
		t.Fatal("expected .png not to match")
	}
	if !matchesExtension("notes.txt", nil) {
		t.Fatal("expected all extensions to match when filter is empty")
	}
}

func TestRunFiltersByExtension(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string][]byte{
		"readme.md": []byte("hello\r\nworld\r\n"),
		"notes.txt": []byte("text\r\nfile\r\n"),
		"data.json": []byte("json\r\nfile\r\n"),
	} {
		writeTestFile(t, root, name, content)
	}

	report, err := Run(t.Context(), Options{Root: root, Extensions: []string{"md", "txt"}, Yes: true})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 2 {
		t.Fatalf("Scanned = %d, want 2", report.Scanned)
	}
	if report.Modified != 2 {
		t.Fatalf("Modified = %d, want 2", report.Modified)
	}
	if report.Skipped != 1 {
		t.Fatalf("Skipped = %d, want 1", report.Skipped)
	}

	got, err := os.ReadFile(filepath.Join(root, "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "json\r\nfile\r\n" {
		t.Fatalf("json file = %q, want unchanged", got)
	}
}
