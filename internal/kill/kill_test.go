package kill

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

type fakeFinder struct {
	processes []Process
	err       error
}

func (f *fakeFinder) FindListeners(ctx context.Context, port int) ([]Process, error) {
	return f.processes, f.err
}

type fakeKiller struct {
	killed []int
	force  []bool
	err    error
}

func (f *fakeKiller) Kill(ctx context.Context, pid int, force bool) error {
	if f.err != nil {
		return f.err
	}
	f.killed = append(f.killed, pid)
	f.force = append(f.force, force)
	return nil
}

func TestRunInvalidPort(t *testing.T) {
	_, err := Run(context.Background(), Options{Port: 0})
	if err == nil || !strings.Contains(err.Error(), "port must be") {
		t.Fatalf("Run() error = %v", err)
	}

	_, err = Run(context.Background(), Options{Port: 70000})
	if err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestRunNoListeners(t *testing.T) {
	result, err := Run(context.Background(), Options{
		Port:   59999,
		Finder: &fakeFinder{},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Found) != 0 {
		t.Fatalf("Found = %v, want empty", result.Found)
	}
}

func TestRunDryRun(t *testing.T) {
	result, err := Run(context.Background(), Options{
		Port:   8080,
		DryRun: true,
		Finder: &fakeFinder{processes: []Process{{PID: 42, Name: "kubectl"}}},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.DryRun || len(result.Killed) != 0 {
		t.Fatalf("result = %+v", result)
	}
	if len(result.Found) != 1 || result.Found[0].PID != 42 {
		t.Fatalf("Found = %+v", result.Found)
	}
}

func TestRunKillsWithYes(t *testing.T) {
	killer := &fakeKiller{}
	result, err := Run(context.Background(), Options{
		Port:   8080,
		Yes:    true,
		Force:  true,
		Finder: &fakeFinder{processes: []Process{{PID: 42, Name: "kubectl"}, {PID: 43, Name: "node"}}},
		Killer: killer,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Killed) != 2 {
		t.Fatalf("Killed = %+v", result.Killed)
	}
	if len(killer.killed) != 2 || killer.killed[0] != 42 || killer.killed[1] != 43 {
		t.Fatalf("killer.killed = %v", killer.killed)
	}
	if len(killer.force) != 2 || !killer.force[0] || !killer.force[1] {
		t.Fatalf("killer.force = %v", killer.force)
	}
}

func TestRunKillError(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Port:   8080,
		Yes:    true,
		Finder: &fakeFinder{processes: []Process{{PID: 42}}},
		Killer: &fakeKiller{err: errors.New("permission denied")},
	})
	if err == nil || !strings.Contains(err.Error(), "kill pid 42") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunNotTerminalRequiresYes(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Port:   8080,
		Finder: &fakeFinder{processes: []Process{{PID: 42}}},
		Input:  strings.NewReader(""),
		Output: io.Discard,
	})
	if err == nil || !strings.Contains(err.Error(), "not a terminal") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestFilterProtected(t *testing.T) {
	filtered := filterProtected([]Process{
		{PID: 1},
		{PID: 42},
		{PID: 42},
	})
	if len(filtered) != 1 || filtered[0].PID != 42 {
		t.Fatalf("filtered = %+v", filtered)
	}
}

func TestRunFinderError(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Port:   8080,
		Finder: &fakeFinder{err: errors.New("finder failed")},
	})
	if err == nil || !strings.Contains(err.Error(), "finder failed") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Run(ctx, Options{
		Port:   8080,
		Finder: contextAwareFinder{},
		Yes:    true,
	})
	if err != context.Canceled {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
}

type contextAwareFinder struct{}

func (contextAwareFinder) FindListeners(ctx context.Context, port int) ([]Process, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func TestRunCancelled(t *testing.T) {
	tty, master := openKillTTY(t)
	defer func() { _ = master.Close() }()
	defer func() { _ = tty.Close() }()

	done := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Millisecond)
		_, _ = master.Write([]byte("n\n"))
		close(done)
	}()

	result, err := Run(context.Background(), Options{
		Port:   8080,
		Finder: &fakeFinder{processes: []Process{{PID: 42}}},
		Input:  tty,
		Output: tty,
	})
	<-done
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Cancelled {
		t.Fatal("expected cancelled result")
	}
}

func TestRunConfirmedKill(t *testing.T) {
	tty, master := openKillTTY(t)
	defer func() { _ = master.Close() }()
	defer func() { _ = tty.Close() }()

	done := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Millisecond)
		_, _ = master.Write([]byte("y\n"))
		close(done)
	}()

	killer := &fakeKiller{}
	result, err := Run(context.Background(), Options{
		Port:   8080,
		Finder: &fakeFinder{processes: []Process{{PID: 42}}},
		Input:  tty,
		Output: tty,
		Killer: killer,
	})
	<-done
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Cancelled || len(result.Killed) != 1 {
		t.Fatalf("result = %+v", result)
	}
}

func openKillTTY(t *testing.T) (*os.File, *os.File) {
	t.Helper()

	master, tty, err := pty.Open()
	if err != nil {
		t.Skip("pty:", err)
	}
	t.Cleanup(func() {
		_ = master.Close()
		_ = tty.Close()
	})
	return tty, master
}

func TestRunKillsWithTermSignal(t *testing.T) {
	killer := &fakeKiller{}
	result, err := Run(context.Background(), Options{
		Port:   8080,
		Yes:    true,
		Force:  false,
		Finder: &fakeFinder{processes: []Process{{PID: 42, Name: "node"}}},
		Killer: killer,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Killed) != 1 {
		t.Fatalf("Killed = %+v", result.Killed)
	}
	if len(killer.force) != 1 || killer.force[0] {
		t.Fatalf("killer.force = %v, want SIGTERM", killer.force)
	}
}

func TestUnsupportedFinder(t *testing.T) {
	_, err := unsupportedFinder{}.FindListeners(context.Background(), 8080)
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("FindListeners() error = %v", err)
	}
}

func TestDefaultFinderAndKiller(t *testing.T) {
	if DefaultFinder() == nil {
		t.Fatal("DefaultFinder() returned nil")
	}
	if DefaultKiller() == nil {
		t.Fatal("DefaultKiller() returned nil")
	}
}

func TestRunConfirmPromptError(t *testing.T) {
	tty, master := openKillTTY(t)

	done := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Millisecond)
		_ = master.Close()
		close(done)
	}()

	_, err := Run(context.Background(), Options{
		Port:   8080,
		Finder: &fakeFinder{processes: []Process{{PID: 42}}},
		Input:  tty,
		Output: tty,
	})
	<-done
	if err == nil {
		t.Fatal("expected confirm prompt error")
	}
}

func TestRunUsesDefaultInputOutput(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Port:   8080,
		Yes:    true,
		Finder: &fakeFinder{},
		Input:  nil,
		Output: nil,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}
