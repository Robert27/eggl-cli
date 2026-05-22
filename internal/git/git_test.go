package git

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	inside      bool
	insideErr   error
	hash        string
	commitErr   error
	pushErr     error
	commitMsg   string
	pushCalled  bool
	commitCalls int
}

func (f *fakeRunner) InsideWorkTree(ctx context.Context) (bool, error) {
	return f.inside, f.insideErr
}

func (f *fakeRunner) EmptyCommit(ctx context.Context, message string) (string, error) {
	f.commitCalls++
	f.commitMsg = message
	if f.commitErr != nil {
		return "", f.commitErr
	}
	return f.hash, nil
}

func (f *fakeRunner) Push(ctx context.Context) error {
	f.pushCalled = true
	return f.pushErr
}

func TestRunNotInsideWorkTree(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Git: &fakeRunner{inside: false},
	})
	if err == nil {
		t.Fatal("expected error when not inside work tree")
	}
	if !strings.Contains(err.Error(), "not inside a git work tree") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunInsideWorkTreeError(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Git: &fakeRunner{insideErr: errors.New("git missing")},
	})
	if err == nil {
		t.Fatal("expected error from InsideWorkTree")
	}
}

func TestRunEmptyCommitOnly(t *testing.T) {
	runner := &fakeRunner{inside: true, hash: "abc1234"}

	result, err := Run(context.Background(), Options{Git: runner})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Hash != "abc1234" {
		t.Fatalf("Hash = %q, want abc1234", result.Hash)
	}
	if result.Pushed {
		t.Fatal("expected Pushed to be false")
	}
	if runner.commitMsg != DefaultMessage {
		t.Fatalf("commit message = %q, want %q", runner.commitMsg, DefaultMessage)
	}
	if runner.pushCalled {
		t.Fatal("Push should not be called")
	}
}

func TestRunWithPush(t *testing.T) {
	runner := &fakeRunner{inside: true, hash: "def5678"}

	result, err := Run(context.Background(), Options{Git: runner, Push: true})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Pushed {
		t.Fatal("expected Pushed to be true")
	}
	if !runner.pushCalled {
		t.Fatal("Push should be called")
	}
}

func TestRunPushError(t *testing.T) {
	runner := &fakeRunner{
		inside:  true,
		hash:    "abc1234",
		pushErr: errors.New("push failed"),
	}

	_, err := Run(context.Background(), Options{Git: runner, Push: true})
	if err == nil {
		t.Fatal("expected push error")
	}
}

func TestRunCustomMessage(t *testing.T) {
	runner := &fakeRunner{inside: true, hash: "abc1234"}

	_, err := Run(context.Background(), Options{
		Git:     runner,
		Message: "custom message",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if runner.commitMsg != "custom message" {
		t.Fatalf("commit message = %q, want custom message", runner.commitMsg)
	}
}

func TestRunIntegration(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	initGitRepo(t, dir)

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	result, err := Run(context.Background(), Options{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Hash == "" {
		t.Fatal("expected non-empty commit hash")
	}
	if len(result.Hash) < 7 {
		t.Fatalf("hash = %q, expected full git hash", result.Hash)
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s in %s: %v: %s", strings.Join(args, " "), filepath.Base(dir), err, strings.TrimSpace(string(out)))
	}
}
