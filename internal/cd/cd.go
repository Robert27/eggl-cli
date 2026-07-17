package cd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/roberteggl/eggl-cli/internal/config"
)

type Options struct {
	ConfigPath string
}

type Entry struct {
	Name string
	Path string
}

func DefaultOptions(configPath string) Options {
	return Options{ConfigPath: configPath}
}

func List(cfg *config.Config) []Entry {
	if cfg == nil || len(cfg.Directories) == 0 {
		return nil
	}

	names := cfg.DirectoryNames()
	sort.Strings(names)

	entries := make([]Entry, 0, len(names))
	for _, name := range names {
		entries = append(entries, Entry{
			Name: name,
			Path: cfg.Directories[name],
		})
	}
	return entries
}

func Resolve(name string, opts Options) (string, error) {
	cfg, path, err := loadConfig(opts.ConfigPath)
	if err != nil {
		return "", err
	}

	dir, ok := cfg.Directories[name]
	if !ok {
		names := cfg.DirectoryNames()
		sort.Strings(names)
		return "", fmt.Errorf("unknown directory %q (config: %s); available: %s",
			name, path, strings.Join(names, ", "))
	}

	expanded, err := config.ExpandPath(dir)
	if err != nil {
		return "", fmt.Errorf("directory %q: %w", name, err)
	}

	abs, err := filepath.Abs(expanded)
	if err != nil {
		return expanded, nil
	}
	return abs, nil
}

func loadConfig(path string) (*config.Config, string, error) {
	if path == "" {
		path = config.DefaultPath()
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, path, err
	}
	return cfg, path, nil
}
