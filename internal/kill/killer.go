package kill

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

type Killer interface {
	Kill(ctx context.Context, pid int, force bool) error
}

func DefaultKiller() Killer {
	if runtime.GOOS == "windows" {
		return windowsKiller{}
	}
	return unixKiller{}
}

type unixKiller struct{}

func (unixKiller) Kill(ctx context.Context, pid int, force bool) error {
	sig := syscall.SIGTERM
	if force {
		sig = syscall.SIGKILL
	}
	if err := syscall.Kill(pid, sig); err != nil {
		return fmt.Errorf("signal %s: %w", sig, err)
	}
	return nil
}

type windowsKiller struct{}

func (windowsKiller) Kill(ctx context.Context, pid int, force bool) error {
	args := []string{"/PID", fmt.Sprintf("%d", pid)}
	if force {
		args = append(args, "/F")
	}
	cmd := exec.CommandContext(ctx, "taskkill", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("taskkill: %w: %s", err, trimOutput(out))
	}
	return nil
}

func processName(pid int) string {
	if runtime.GOOS == "windows" {
		return ""
	}

	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}
	return trimOutput(data)
}
