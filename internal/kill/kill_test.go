package kill

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
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

func TestParseProcNetTCP(t *testing.T) {
	data := `  sl  local_address rem_address   st tx queue rx queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 4242 1 00000000ab12cd34 0 0 0 0 -1
   1: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 9999 1 00000000ab12cd35 0 0 0 0 -1
   2: 0100007F:1F91 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 1111 1 00000000ab12cd36 0 0 0 0 -1`

	inodes, err := parseProcNetTCP(data, 8080)
	if err != nil {
		t.Fatalf("parseProcNetTCP() error = %v", err)
	}
	if len(inodes) != 2 || inodes[0] != 4242 || inodes[1] != 9999 {
		t.Fatalf("inodes = %v", inodes)
	}

	inodes, err = parseProcNetTCP(data, 8081)
	if err != nil {
		t.Fatalf("parseProcNetTCP() error = %v", err)
	}
	if len(inodes) != 1 || inodes[0] != 1111 {
		t.Fatalf("inodes = %v", inodes)
	}
}
