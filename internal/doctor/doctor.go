package doctor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
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

	checks := make([]Check, 0, 3)

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

	return &Report{Checks: checks}, nil
}

func HasFailures(report *Report) bool {
	for _, check := range report.Checks {
		if !check.OK {
			return true
		}
	}
	return false
}
