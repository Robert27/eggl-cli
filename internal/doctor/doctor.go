package doctor

import (
	"context"
	"fmt"
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

	checks := []Check{
		{
			Name:   "go",
			Status: runtime.Version(),
			Detail: "Go runtime available",
			OK:     true,
		},
		{
			Name:   "os",
			Status: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			Detail: "Platform detected",
			OK:     true,
		},
	}

	path := opts.CheckPath
	if path == "" {
		path = os.Getenv("HOME")
	}

	if path == "" {
		checks = append(checks, Check{
			Name:   "home",
			Status: "missing",
			Detail: "HOME environment variable is not set",
			OK:     false,
		})
	} else if info, err := os.Stat(path); err != nil {
		checks = append(checks, Check{
			Name:   "home",
			Status: "error",
			Detail: err.Error(),
			OK:     false,
		})
	} else if !info.IsDir() {
		checks = append(checks, Check{
			Name:   "home",
			Status: "invalid",
			Detail: fmt.Sprintf("%s is not a directory", path),
			OK:     false,
		})
	} else {
		checks = append(checks, Check{
			Name:   "home",
			Status: path,
			Detail: "Home directory accessible",
			OK:     true,
		})
	}

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
