package dedash

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(root, name, content string) error {
	return os.WriteFile(filepath.Join(root, name), []byte(content), 0o644)
}

func readTestFile(root, name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestRunReplacesEmDash(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	if err := os.WriteFile(path, []byte("hello \u2014 world \u2014 again"), 0o644); err != nil {
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
	original := "hello \u2014 world"
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
	if err := os.WriteFile(path, []byte("fake\x00png - dash"), 0o644); err != nil {
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
	if err := os.WriteFile(path, []byte("hello \u2014 world"), 0o644); err != nil {
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
	if string(got) != "hello \u2014 world" {
		t.Fatalf("file content = %q, want unchanged", got)
	}
}

func TestRunDoesNotReplaceEnDash(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "readme.md")
	original := "range \u2013 value"
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

func TestShouldSkipFile(t *testing.T) {
	if !shouldSkipFile("photo.JPG") {
		t.Fatal("expected .jpg to be skipped")
	}
	if shouldSkipFile("readme.md") {
		t.Fatal("expected .md not to be skipped")
	}
	if !shouldSkipFile("app.jar") {
		t.Fatal("expected .jar to be skipped")
	}
	if !shouldSkipFile("bundle.min.js") {
		t.Fatal("expected .min.js to be skipped")
	}
	if !shouldSkipFile(".DS_Store") {
		t.Fatal("expected .DS_Store to be skipped")
	}
}

func TestRunSkipsBuildOutputDirs(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		dir  string
		file string
	}{
		{"target", "classes/readme.md"},
		{".gradle", "caches/readme.md"},
		{"__pycache__", "module/readme.md"},
	}
	for _, tc := range cases {
		path := filepath.Join(root, tc.dir, tc.file)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("hello \u2014 world"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 0 {
		t.Fatalf("Scanned = %d, want 0", report.Scanned)
	}
	if report.Modified != 0 {
		t.Fatalf("Modified = %d, want 0", report.Modified)
	}
}

func TestRunSkipsJarFiles(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "lib", "app.jar")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("PK\x03\x04fake jar"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 0 {
		t.Fatalf("Scanned = %d, want 0", report.Scanned)
	}
	if report.Skipped != 1 {
		t.Fatalf("Skipped = %d, want 1", report.Skipped)
	}
}

func TestIsBinaryContent(t *testing.T) {
	if !isBinaryContent([]byte("hello\x00world")) {
		t.Fatal("expected null byte content to be binary")
	}
	if isBinaryContent([]byte("hello \u2014 world")) {
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

func TestRunScannedWithoutEmDash(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "plain.md")
	if err := os.WriteFile(path, []byte("no special dashes here"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 1 {
		t.Fatalf("Scanned = %d, want 1", report.Scanned)
	}
	if report.Modified != 0 {
		t.Fatalf("Modified = %d, want 0", report.Modified)
	}
}

func TestRunMultipleFiles(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string]string{
		"a.md": "one \u2014 dash",
		"b.md": "two \u2014 dashes \u2014 here",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 2 {
		t.Fatalf("Modified = %d, want 2", report.Modified)
	}
	if TotalReplacements(report) != 3 {
		t.Fatalf("TotalReplacements = %d, want 3", TotalReplacements(report))
	}
}

func TestRunSkipsNonRegularFile(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "readme.md")
	if err := os.WriteFile(target, []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link.md")
	if err := os.Symlink(target, link); err != nil {
		t.Skip("symlinks not supported:", err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 1 {
		t.Fatalf("Scanned = %d, want 1 (regular file only)", report.Scanned)
	}
	if report.Skipped != 1 {
		t.Fatalf("Skipped = %d, want 1 (symlink)", report.Skipped)
	}
}

func TestIsBinaryContentEmpty(t *testing.T) {
	if isBinaryContent(nil) {
		t.Fatal("expected empty content not to be binary")
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

func TestRunSkipsHiddenFile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("secret \u2014 value"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "visible.md"), []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 1 {
		t.Fatalf("Scanned = %d, want 1", report.Scanned)
	}
	if report.Skipped != 1 {
		t.Fatalf("Skipped = %d, want 1 (.env)", report.Skipped)
	}
}

func TestRunIncludeHidden(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("secret \u2014 value"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root, IncludeHidden: true})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Modified != 1 {
		t.Fatalf("Modified = %d, want 1", report.Modified)
	}
}

func TestRunSkipsOversizedFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "big.md")
	if err := os.WriteFile(path, []byte("\u2014"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(path, MaxFileSize+1); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.Scanned != 0 {
		t.Fatalf("Scanned = %d, want 0", report.Scanned)
	}
	if report.Skipped != 1 {
		t.Fatalf("Skipped = %d, want 1", report.Skipped)
	}
}

func TestRunRespectsContextCancel(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "readme.md"), []byte("hello \u2014 world"), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Run(ctx, Options{Root: root})
	if err != context.Canceled {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
}
