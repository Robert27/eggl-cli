package pf

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

func defaultOpenURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open %s: %w", url, err)
	}
	return nil
}

func waitForLocalPort(ctx context.Context, port string) error {
	addr := net.JoinHostPort("127.0.0.1", port)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func waitUntilReady(ctx context.Context, port string, pfErr <-chan error, ready func(context.Context, string) error) error {
	readyCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	readyErr := make(chan error, 1)
	go func() {
		readyErr <- ready(readyCtx, port)
	}()

	for {
		select {
		case err := <-pfErr:
			cancel()
			if err != nil {
				return err
			}
			return fmt.Errorf("port-forward closed before localhost:%s was ready", port)
		case err := <-readyErr:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
