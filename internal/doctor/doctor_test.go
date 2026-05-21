package doctor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunUsesCheckPath(t *testing.T) {
	dir := t.TempDir()

	report, err := Run(context.Background(), Options{CheckPath: dir})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if !home.OK {
		t.Fatalf("home check should pass, got %+v", home)
	}
	if home.Status != dir {
		t.Fatalf("home Status = %q, want %q", home.Status, dir)
	}
}

func TestRunDefaultUsesHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	report, err := Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if !home.OK {
		t.Fatalf("home check should pass, got %+v", home)
	}
	if home.Status != dir {
		t.Fatalf("home Status = %q, want %q", home.Status, dir)
	}
}

func TestRunMissingHome(t *testing.T) {
	t.Setenv("HOME", "")

	report, err := Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if home.OK {
		t.Fatal("expected home check to fail when HOME is unset")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunInvalidCheckPath(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")

	report, err := Run(context.Background(), Options{CheckPath: missing})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if home.OK {
		t.Fatal("expected home check to fail for missing path")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunCheckPathNotDirectory(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	if err := os.WriteFile(path, []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(context.Background(), Options{CheckPath: path})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	home := findCheck(t, report, "home")
	if home.OK {
		t.Fatal("expected home check to fail for non-directory path")
	}
	if !HasFailures(report) {
		t.Fatal("expected report to have failures")
	}
}

func TestRunIncludesRuntimeChecks(t *testing.T) {
	dir := t.TempDir()

	report, err := Run(context.Background(), Options{CheckPath: dir})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for _, name := range []string{"go", "os"} {
		check := findCheck(t, report, name)
		if !check.OK {
			t.Fatalf("%s check should pass, got %+v", name, check)
		}
	}
}

func TestHasFailures(t *testing.T) {
	if HasFailures(&Report{Checks: []Check{{OK: true}}}) {
		t.Fatal("expected no failures")
	}
	if !HasFailures(&Report{Checks: []Check{{OK: false}}}) {
		t.Fatal("expected failures")
	}
}

func findCheck(t *testing.T, report *Report, name string) Check {
	t.Helper()

	for _, check := range report.Checks {
		if check.Name == name {
			return check
		}
	}
	t.Fatalf("check %q not found in report", name)
	return Check{}
}
