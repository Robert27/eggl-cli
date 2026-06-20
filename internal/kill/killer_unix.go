//go:build !windows

package kill

import (
	"context"
	"fmt"
	"os"
	"syscall"
)

type unixKiller struct{}

func platformKiller() Killer {
	return unixKiller{}
}

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

func processName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}
	return trimOutput(data)
}
