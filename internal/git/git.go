package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const DefaultMessage = "chore: empty commit"

type Runner interface {
	InsideWorkTree(ctx context.Context) (bool, error)
	EmptyCommit(ctx context.Context, message string) (string, error)
	Push(ctx context.Context) error
}

type CLI struct {
	Bin string
}

func (c CLI) bin() string {
	if c.Bin != "" {
		return c.Bin
	}
	return "git"
}

func (c CLI) InsideWorkTree(ctx context.Context) (bool, error) {
	out, err := exec.CommandContext(ctx, c.bin(), "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false, fmt.Errorf("git rev-parse --is-inside-work-tree: %w", err)
	}
	return strings.TrimSpace(string(out)) == "true", nil
}

func (c CLI) EmptyCommit(ctx context.Context, message string) (string, error) {
	cmd := exec.CommandContext(ctx, c.bin(), "commit", "--allow-empty", "-m", message)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git commit --allow-empty: %w: %s", err, strings.TrimSpace(string(out)))
	}

	hash, err := exec.CommandContext(ctx, c.bin(), "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse HEAD: %w", err)
	}
	return strings.TrimSpace(string(hash)), nil
}

func (c CLI) Push(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.bin(), "push")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

type Options struct {
	Message string
	Push    bool
	Git     Runner
}

type Result struct {
	Hash   string
	Pushed bool
}

func Run(ctx context.Context, opts Options) (*Result, error) {
	runner := opts.Git
	if runner == nil {
		runner = CLI{}
	}

	message := opts.Message
	if message == "" {
		message = DefaultMessage
	}

	inside, err := runner.InsideWorkTree(ctx)
	if err != nil {
		return nil, err
	}
	if !inside {
		return nil, fmt.Errorf("not inside a git work tree")
	}

	hash, err := runner.EmptyCommit(ctx, message)
	if err != nil {
		return nil, err
	}

	result := &Result{Hash: hash}
	if opts.Push {
		if err := runner.Push(ctx); err != nil {
			return nil, err
		}
		result.Pushed = true
	}

	return result, nil
}
