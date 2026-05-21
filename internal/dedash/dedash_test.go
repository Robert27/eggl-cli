package dedash

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReplacesEmDash(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello — world — again"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if report.Modified != 1 {
		t.Fatalf("Modified = %d, want 1", report.Modified)
	}
	if TotalReplacements(report) != 2 {
		t.Fatalf("TotalReplacements = %d, want 2", TotalReplacements(report))
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello - world - again" {
		t.Fatalf("file content = %q, want %q", got, "hello - world - again")
	}
}

func TestRunDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	original := "hello — world"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root, DryRun: true})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if report.Modified != 1 {
		t.Fatalf("Modified = %d, want 1", report.Modified)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != original {
		t.Fatalf("file content = %q, want unchanged %q", got, original)
	}
}

func TestRunSkipsBinaryFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "image.png")
	if err := os.WriteFile(path, []byte("fake\x00png — dash"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if report.Modified != 0 {
		t.Fatalf("Modified = %d, want 0", report.Modified)
	}
	if report.Skipped != 1 {
		t.Fatalf("Skipped = %d, want 1", report.Skipped)
	}
}

func TestRunSkipsNodeModules(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "node_modules", "pkg", "index.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("hello — world"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if report.Modified != 0 {
		t.Fatalf("Modified = %d, want 0", report.Modified)
	}
	if report.Scanned != 0 {
		t.Fatalf("Scanned = %d, want 0", report.Scanned)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello — world" {
		t.Fatalf("file content = %q, want unchanged", got)
	}
}

func TestRunDoesNotReplaceEnDash(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	original := "range – value"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if report.Modified != 0 {
		t.Fatalf("Modified = %d, want 0", report.Modified)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != original {
		t.Fatalf("file content = %q, want unchanged %q", got, original)
	}
}

func TestShouldSkipExtension(t *testing.T) {
	if !shouldSkipExtension("photo.JPG") {
		t.Fatal("expected .jpg to be skipped")
	}
	if shouldSkipExtension("readme.md") {
		t.Fatal("expected .md not to be skipped")
	}
}

func TestIsBinaryContent(t *testing.T) {
	if !isBinaryContent([]byte("hello\x00world")) {
		t.Fatal("expected null byte content to be binary")
	}
	if isBinaryContent([]byte("hello — world")) {
		t.Fatal("expected valid UTF-8 text not to be binary")
	}
	if !isBinaryContent([]byte{0xff, 0xfe, 0xfd}) {
		t.Fatal("expected invalid UTF-8 to be treated as binary")
	}
}

func TestRunInvalidPath(t *testing.T) {
	_, err := Run(context.Background(), Options{Root: filepath.Join(t.TempDir(), "missing")})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRunFileNotDirectory(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	if err := os.WriteFile(path, []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Run(context.Background(), Options{Root: path})
	if err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("expected not a directory error, got %v", err)
	}
}
