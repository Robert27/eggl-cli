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
