package netbird

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type Profile struct {
	Name   string
	Active bool
}

type Runner interface {
	ListProfiles(ctx context.Context) ([]Profile, error)
	SelectProfile(ctx context.Context, name string) error
}

type CLI struct {
	Bin string
}

func (c CLI) bin() string {
	if c.Bin != "" {
		return c.Bin
	}
	return "netbird"
}

func (c CLI) ListProfiles(ctx context.Context) ([]Profile, error) {
	cmd := exec.CommandContext(ctx, c.bin(), "profile", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("netbird profile list: %w", err)
	}
	return ParseProfileList(out)
}

func (c CLI) SelectProfile(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, c.bin(), "profile", "select", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("netbird profile select %s: %w: %s", name, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func ParseProfileList(data []byte) ([]Profile, error) {
	var profiles []Profile
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Found ") {
			continue
		}
		active := false
		name := line
		switch {
		case strings.HasPrefix(line, "✓"):
			active = true
			name = strings.TrimSpace(strings.TrimPrefix(line, "✓"))
		case strings.HasPrefix(line, "✗"):
			name = strings.TrimSpace(strings.TrimPrefix(line, "✗"))
		}
		if name == "" {
			continue
		}
		profiles = append(profiles, Profile{Name: name, Active: active})
	}
	return profiles, nil
}

func ResolveProfile(ref string, profiles []Profile) (Profile, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Profile{}, fmt.Errorf("empty netbird profile reference")
	}

	var matches []Profile
	for _, p := range profiles {
		if strings.EqualFold(p.Name, ref) {
			matches = append(matches, p)
		}
	}

	switch len(matches) {
	case 0:
		return Profile{}, fmt.Errorf("netbird profile %q not found in profile list", ref)
	case 1:
		return matches[0], nil
	default:
		names := make([]string, len(matches))
		for i, p := range matches {
			names[i] = p.Name
		}
		sort.Strings(names)
		return Profile{}, fmt.Errorf("ambiguous netbird profile reference %q: matches %s", ref, strings.Join(names, ", "))
	}
}

func CurrentProfile(profiles []Profile) (Profile, error) {
	for _, p := range profiles {
		if p.Active {
			return p, nil
		}
	}
	return Profile{}, fmt.Errorf("no active netbird profile")
}

func FormatProfile(p Profile) string {
	return p.Name
}
