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

	"github.com/creack/pty"
	"github.com/roberteggl/eggl-cli/internal/kill"
)

func TestMain(m *testing.M) {
	portStr := os.Getenv("EGGL_HOLD_PORT")
	if portStr == "" {
		os.Exit(m.Run())
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EGGL_HOLD_PORT: %v\n", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

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

	cmd := exec.Command(os.Args[0], "-test.paniconexit0")
	cmd.Env = append(os.Environ(), "EGGL_HOLD_PORT="+strconv.Itoa(port))
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

func TestKillYesWithoutForce(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping listener test in short mode")
	}

	t.Setenv("NO_COLOR", "1")

	port := startListenerProcess(t)
	stdout, _, err := runCmd(t, "kill", "--yes", strconv.Itoa(port))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stdout, "killed pid") {
		t.Fatalf("stdout = %q", stdout)
	}
	if strings.Contains(stdout, "SIGKILL") {
		t.Fatalf("expected SIGTERM kill output, got %q", stdout)
	}
}

func TestKillCancelled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping listener test in short mode")
	}

	t.Setenv("NO_COLOR", "1")

	port := startListenerProcess(t)
	tty, master := openKillCmdTTY(t)
	defer func() { _ = master.Close() }()
	defer func() { _ = tty.Close() }()

	done := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Millisecond)
		_, _ = master.Write([]byte("n\n"))
		close(done)
	}()

	stdout, stderr, err := runCmdWithIn(t, tty, "kill", strconv.Itoa(port))
	<-done
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(stderr, "cancelled") {
		t.Fatalf("stderr = %q", stderr)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
}

func openKillCmdTTY(t *testing.T) (*os.File, *os.File) {
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
