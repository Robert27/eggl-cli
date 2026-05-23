package kube

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Runner interface {
	CurrentContext(ctx context.Context) (string, error)
	UseContext(ctx context.Context, name string) error
	PortForward(ctx context.Context, args []string) error
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

func (c CLI) PortForward(ctx context.Context, args []string) error {
	kubectlArgs := append([]string{"port-forward"}, args...)
	cmd := exec.CommandContext(ctx, c.bin(), kubectlArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl port-forward: %w", err)
	}
	return nil
}
