package doctor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"

	"github.com/roberteggl/eggl-cli/internal/config"
)

type Options struct {
	CheckPath string
}

type Check struct {
	Name   string
	Status string
	Detail string
	OK     bool
}

type Report struct {
	Checks []Check
}

func Run(ctx context.Context, opts Options) (*Report, error) {
	_ = ctx

	slog.Debug("starting environment checks", "check_path", opts.CheckPath)

	checks := make([]Check, 0, 7)

	slog.Debug("running check", "name", "go")
	checks = append(checks, Check{
		Name:   "go",
		Status: runtime.Version(),
		Detail: "Go runtime available",
		OK:     true,
	})

	slog.Debug("running check", "name", "os")
	checks = append(checks, Check{
		Name:   "os",
		Status: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Detail: "Platform detected",
		OK:     true,
	})

	path := opts.CheckPath
	if path == "" {
		path = os.Getenv("HOME")
	}

	slog.Debug("running check", "name", "home", "path", path)

	var homeCheck Check
	if path == "" {
		homeCheck = Check{
			Name:   "home",
			Status: "missing",
			Detail: "HOME environment variable is not set",
			OK:     false,
		}
	} else if info, err := os.Stat(path); err != nil {
		homeCheck = Check{
			Name:   "home",
			Status: "error",
			Detail: err.Error(),
			OK:     false,
		}
	} else if !info.IsDir() {
		homeCheck = Check{
			Name:   "home",
			Status: "invalid",
			Detail: fmt.Sprintf("%s is not a directory", path),
			OK:     false,
		}
	} else {
		homeCheck = Check{
			Name:   "home",
			Status: path,
			Detail: "Home directory accessible",
			OK:     true,
		}
	}
	checks = append(checks, homeCheck)
	slog.Debug("check result",
		"name", homeCheck.Name,
		"ok", homeCheck.OK,
		"status", homeCheck.Status,
		"detail", homeCheck.Detail,
	)

	for _, tool := range []string{"kubectl", "git", "tailscale"} {
		check := toolCheck(tool)
		checks = append(checks, check)
		slog.Debug("check result",
			"name", check.Name,
			"ok", check.OK,
			"status", check.Status,
			"detail", check.Detail,
		)
	}

	configCheck := configCheck()
	checks = append(checks, configCheck)
	slog.Debug("check result",
		"name", configCheck.Name,
		"ok", configCheck.OK,
		"status", configCheck.Status,
		"detail", configCheck.Detail,
	)

	return &Report{Checks: checks}, nil
}

func toolCheck(name string) Check {
	path, err := exec.LookPath(name)
	if err != nil {
		return Check{
			Name:   name,
			Status: "missing",
			Detail: fmt.Sprintf("%s not found on PATH", name),
			OK:     false,
		}
	}
	return Check{
		Name:   name,
		Status: path,
		Detail: fmt.Sprintf("%s available on PATH", name),
		OK:     true,
	}
}

func configCheck() Check {
	path := config.DefaultPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Check{
			Name:   "config",
			Status: "not found",
			Detail: "run eggl env init",
			OK:     true,
		}
	} else if err != nil {
		return Check{
			Name:   "config",
			Status: "error",
			Detail: err.Error(),
			OK:     false,
		}
	}

	if _, err := config.Load(path); err != nil {
		return Check{
			Name:   "config",
			Status: "invalid",
			Detail: err.Error(),
			OK:     false,
		}
	}

	return Check{
		Name:   "config",
		Status: path,
		Detail: "Config file valid",
		OK:     true,
	}
}

func HasFailures(report *Report) bool {
	for _, check := range report.Checks {
		if !check.OK {
			return true
		}
	}
	return false
}
