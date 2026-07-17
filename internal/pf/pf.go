package pf

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/roberteggl/eggl-cli/internal/config"
	"github.com/roberteggl/eggl-cli/internal/kube"
)

const defaultPorts = "8080:80"

type Options struct {
	ConfigPath string
	Kube       kube.Runner
}

type Entry struct {
	Name      string
	Namespace string
	Resource  string
	Ports     []string
}

func DefaultOptions(configPath string) Options {
	return Options{
		ConfigPath: configPath,
		Kube:       kube.CLI{},
	}
}

func List(cfg *config.Config) []Entry {
	if cfg == nil || len(cfg.PortForwards) == 0 {
		return nil
	}

	names := cfg.PortForwardNames()
	sort.Strings(names)

	entries := make([]Entry, 0, len(names))
	for _, name := range names {
		pf := cfg.PortForwards[name]
		entries = append(entries, Entry{
			Name:      name,
			Namespace: pf.Namespace,
			Resource:  pf.Resource,
			Ports:     resolvedPorts(pf.Ports),
		})
	}
	return entries
}

func Run(ctx context.Context, name string, opts Options) error {
	cfg, path, err := loadConfig(opts.ConfigPath)
	if err != nil {
		return err
	}

	pf, ok := cfg.PortForwards[name]
	if !ok {
		names := cfg.PortForwardNames()
		sort.Strings(names)
		return fmt.Errorf("unknown port-forward %q (config: %s); available: %s",
			name, path, strings.Join(names, ", "))
	}

	ports := resolvedPorts(pf.Ports)
	args := append([]string{"-n", pf.Namespace, pf.Resource}, ports...)

	localPort := strings.Split(ports[0], ":")[0]
	fmt.Fprintf(os.Stderr, "port-forward %s → localhost:%s (%s/%s)\n",
		name, localPort, pf.Namespace, pf.Resource)

	if opts.Kube == nil {
		opts.Kube = kube.CLI{}
	}
	return opts.Kube.PortForward(ctx, args)
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

func resolvedPorts(ports []string) []string {
	if len(ports) == 0 {
		return []string{defaultPorts}
	}
	return ports
}
