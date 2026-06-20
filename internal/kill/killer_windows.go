//go:build windows

package kill

import (
	"context"
	"fmt"
	"os/exec"
)

type windowsKiller struct{}

func platformKiller() Killer {
	return windowsKiller{}
}

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

func processName(int) string {
	return ""
}
