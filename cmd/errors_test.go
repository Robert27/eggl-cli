package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandErrors(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "dedash missing path",
			args:    []string{"dedash", "--path", filepath.Join(t.TempDir(), "missing")},
			wantErr: "no such file or directory",
		},
		{
			name:    "dedash path not directory",
			args:    []string{"dedash", "--path", writeTempFile(t)},
			wantErr: "not a directory",
		},
		{
			name:    "completion invalid shell",
			args:    []string{"completion", "nushell"},
			wantErr: "invalid argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := runCmd(t, tt.args...)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.wantErr)) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func writeTempFile(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
