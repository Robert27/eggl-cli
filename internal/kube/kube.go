package kube

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	CurrentContext(ctx context.Context) (string, error)
	UseContext(ctx context.Context, name string) error
}

type CLI struct {
	Bin string
}

func (c CLI) bin() string {
	if c.Bin != "" {
		return c.Bin
	}
	return "kubectl"
}

func (c CLI) CurrentContext(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, c.bin(), "config", "current-context").Output()
	if err != nil {
		return "", fmt.Errorf("kubectl current-context: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (c CLI) UseContext(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, c.bin(), "config", "use-context", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("kubectl use-context %s: %w: %s", name, err, strings.TrimSpace(string(out)))
	}
	return nil
}
