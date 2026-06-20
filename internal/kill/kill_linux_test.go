//go:build linux

package kill

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestHoldPortHelper(t *testing.T) {
	if os.Getenv("KILL_TEST_PORT") == "" {
		t.Skip("helper only")
	}

	port, err := strconv.Atoi(os.Getenv("KILL_TEST_PORT"))
	if err != nil {
		t.Fatalf("KILL_TEST_PORT: %v", err)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	select {}
}

func startListenerProcess(t *testing.T) int {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestHoldPortHelper$", "-test.count=1")
	cmd.Env = append(os.Environ(), "KILL_TEST_PORT="+strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		processes, err := linuxFinder{}.FindListeners(context.Background(), port)
		if err == nil && len(processes) == 1 {
			return port
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatal("port holder did not start")
	return 0
}

func TestLinuxFinderFindListenersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	port := startListenerProcess(t)
	processes, err := linuxFinder{}.FindListeners(context.Background(), port)
	if err != nil {
		t.Fatalf("FindListeners() error = %v", err)
	}
	if len(processes) != 1 {
		t.Fatalf("processes = %+v, want one listener", processes)
	}
}

func TestRunKillsListenerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	port := startListenerProcess(t)
	result, err := Run(context.Background(), Options{
		Port:  port,
		Yes:   true,
		Force: true,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(result.Killed) != 1 {
		t.Fatalf("Killed = %+v", result.Killed)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		processes, err := linuxFinder{}.FindListeners(context.Background(), port)
		if err != nil {
			t.Fatalf("FindListeners() error = %v", err)
		}
		if len(processes) == 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("listener still present on port %d after kill", port)
}

func TestPidsForInodeNoMatch(t *testing.T) {
	pids, err := pidsForInode(999999999)
	if err != nil {
		t.Fatalf("pidsForInode() error = %v", err)
	}
	if len(pids) != 0 {
		t.Fatalf("pids = %v, want none", pids)
	}
}

func TestProcessNameCurrentProcess(t *testing.T) {
	name := processName(os.Getpid())
	if name == "" {
		t.Fatal("expected process name for current pid")
	}
}

func TestProcessNameInvalidPID(t *testing.T) {
	if processName(999999999) != "" {
		t.Fatal("expected empty name for invalid pid")
	}
}

func TestUnixKillerInvalidPID(t *testing.T) {
	err := unixKiller{}.Kill(context.Background(), 99999999, true)
	if err == nil {
		t.Fatal("expected error killing invalid pid")
	}
}

func TestListenerInodesNoListeners(t *testing.T) {
	inodes, err := listenerInodes(59999)
	if err != nil {
		t.Fatalf("listenerInodes() error = %v", err)
	}
	if len(inodes) != 0 {
		t.Fatalf("inodes = %v, want empty", inodes)
	}
}

func TestLinuxFinderContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := linuxFinder{}.FindListeners(ctx, 8080)
	if err != context.Canceled {
		t.Fatalf("FindListeners() error = %v, want context.Canceled", err)
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

func TestParseProcNetTCPSkipsNonListenState(t *testing.T) {
	data := `  sl  local_address rem_address   st tx queue rx queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 0100007F:1F91 01 00000000:00000000 00:00000000 00000000  1000        0 4242 1 00000000ab12cd34 0 0 0 0 -1`

	inodes, err := parseProcNetTCP(data, 8080)
	if err != nil {
		t.Fatalf("parseProcNetTCP() error = %v", err)
	}
	if len(inodes) != 0 {
		t.Fatalf("inodes = %v, want none for non-listen state", inodes)
	}
}

func TestParseProcNetTCPInvalidInode(t *testing.T) {
	data := `  sl  local_address rem_address   st tx queue rx queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 not-an-inode 1 00000000ab12cd34 0 0 0 0 -1`

	_, err := parseProcNetTCP(data, 8080)
	if err == nil {
		t.Fatal("expected parse error for invalid inode")
	}
}

func TestParseProcNetTCPSkipsMalformedLines(t *testing.T) {
	data := `  sl  local_address rem_address   st tx queue rx queue tr tm->when retrnsmt   uid  timeout inode
   0: malformed line without enough fields
   1: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 4242 1 00000000ab12cd34 0 0 0 0 -1
   2: 0100007F:badlocal 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 1111 1 00000000ab12cd36 0 0 0 0 -1`

	inodes, err := parseProcNetTCP(data, 8080)
	if err != nil {
		t.Fatalf("parseProcNetTCP() error = %v", err)
	}
	if len(inodes) != 1 || inodes[0] != 4242 {
		t.Fatalf("inodes = %v", inodes)
	}
}
