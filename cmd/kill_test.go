package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Robert27/eggl-cli/internal/kill"
)

func TestHoldPortHelper(t *testing.T) {
	if os.Getenv("EGGL_KILL_TEST_PORT") == "" {
		t.Skip("helper only")
	}

	port, err := strconv.Atoi(os.Getenv("EGGL_KILL_TEST_PORT"))
	if err != nil {
		t.Fatalf("EGGL_KILL_TEST_PORT: %v", err)
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
	cmd.Env = append(os.Environ(), "EGGL_KILL_TEST_PORT="+strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		result, err := kill.Run(context.Background(), kill.Options{
			Port:   port,
			DryRun: true,
			Yes:    true,
		})
		if err == nil && len(result.Found) == 1 {
			return port
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatal("port holder did not start")
	return 0
}

func TestKillHelp(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := runHelp(t, "kill", "--help")

	for _, want := range []string{
		"eggl kill",
		"--dry-run",
		"--yes",
		"--force",
		"Examples:",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in kill help, got %q", want, output)
		}
	}

	if strings.Contains(output, "Available Commands:") {
		t.Fatalf("kill help should not list subcommands, got %q", output)
	}
}

func TestKillInvalidPort(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	_, _, err := runCmd(t, "kill", "abc")
	if err == nil || !strings.Contains(err.Error(), "invalid port") {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestKillDryRunNoListeners(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stdout, _, err := runCmd(t, "kill", "--dry-run", "59999")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "no process listening on port 59999") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestKillDryRunWithListener(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping listener test in short mode")
	}

	t.Setenv("NO_COLOR", "1")

	port := startListenerProcess(t)
	stdout, _, err := runCmd(t, "kill", "--dry-run", strconv.Itoa(port))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "would kill pid") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestKillYesForce(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping listener test in short mode")
	}

	t.Setenv("NO_COLOR", "1")

	port := startListenerProcess(t)
	stdout, _, err := runCmd(t, "kill", "--yes", "--force", strconv.Itoa(port))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "killed pid") || !strings.Contains(stdout, "SIGKILL") {
		t.Fatalf("stdout = %q", stdout)
	}
}
