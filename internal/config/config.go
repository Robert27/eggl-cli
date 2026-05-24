package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

var (
	portMappingRE = regexp.MustCompile(`^\d+:\d+$`)
	namespaceRE   = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	resourceRE    = regexp.MustCompile(`^(?i)(svc|service|deploy|deployment|pod|sts|statefulset|daemonset|ds|job|cronjob)/[a-z0-9]([-a-z0-9.]*[a-z0-9])?$`)
)

const (
	fileName      = "config.yaml"
	MaxConfigSize = 1 << 20

	VPNTailscale = "tailscale"
	VPNNetbird   = "netbird"
)

type Config struct {
	Profiles     map[string]Profile     `yaml:"profiles"`
	PortForwards map[string]PortForward `yaml:"port_forwards"`
}

type Profile struct {
	KubeContext      string `yaml:"kube_context"`
	VPN              string `yaml:"vpn,omitempty"`
	TailscaleAccount string `yaml:"tailscale_account,omitempty"`
	NetbirdProfile   string `yaml:"netbird_profile,omitempty"`
}

func (p Profile) VPNType() string {
	switch strings.ToLower(strings.TrimSpace(p.VPN)) {
	case "", VPNTailscale:
		return VPNTailscale
	default:
		return strings.ToLower(strings.TrimSpace(p.VPN))
	}
}

type PortForward struct {
	Namespace string   `yaml:"namespace"`
	Resource  string   `yaml:"resource"`
	Ports     []string `yaml:"ports"`
}

func DefaultPath() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "eggl", fileName)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".config", "eggl", fileName)
	}
	return filepath.Join(home, ".config", "eggl", fileName)
}

func Load(path string) (*Config, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	if info.Size() > MaxConfigSize {
		return nil, fmt.Errorf("config: file exceeds maximum size of %d bytes", MaxConfigSize)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Profiles) == 0 {
		return fmt.Errorf("config: at least one profile is required")
	}

	seenTargets := make(map[string]string, len(c.Profiles))
	for name, profile := range c.Profiles {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("config: profile name must not be empty")
		}
		if strings.TrimSpace(profile.KubeContext) == "" {
			return fmt.Errorf("config: profile %q: kube_context is required", name)
		}
		if err := validateKubeContext(name, profile.KubeContext); err != nil {
			return err
		}
		if err := validateProfileVPN(name, profile); err != nil {
			return err
		}

		key := profileTargetKey(profile)
		if other, ok := seenTargets[key]; ok {
			return fmt.Errorf("config: profiles %q and %q share the same kube_context and mesh VPN target", other, name)
		}
		seenTargets[key] = name
	}

	for name, pf := range c.PortForwards {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("config: port_forward name must not be empty")
		}
		if strings.TrimSpace(pf.Namespace) == "" {
			return fmt.Errorf("config: port_forward %q: namespace is required", name)
		}
		if err := validateNamespace(name, pf.Namespace); err != nil {
			return err
		}
		if strings.TrimSpace(pf.Resource) == "" {
			return fmt.Errorf("config: port_forward %q: resource is required", name)
		}
		if err := validateResource(name, pf.Resource); err != nil {
			return err
		}
		for _, port := range pf.Ports {
			if err := validatePortMapping(name, port); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateKubeContext(profileName, context string) error {
	context = strings.TrimSpace(context)
	if strings.HasPrefix(context, "-") {
		return fmt.Errorf("config: profile %q: kube_context must not start with '-'", profileName)
	}
	if strings.ContainsFunc(context, unicode.IsSpace) {
		return fmt.Errorf("config: profile %q: kube_context must not contain whitespace", profileName)
	}
	if len(context) > 253 {
		return fmt.Errorf("config: profile %q: kube_context exceeds maximum length of 253", profileName)
	}
	return nil
}

func validateNamespace(pfName, namespace string) error {
	namespace = strings.TrimSpace(namespace)
	if len(namespace) > 63 {
		return fmt.Errorf("config: port_forward %q: namespace exceeds maximum length of 63", pfName)
	}
	if !namespaceRE.MatchString(namespace) {
		return fmt.Errorf("config: port_forward %q: namespace %q is invalid", pfName, namespace)
	}
	return nil
}

func validateResource(pfName, resource string) error {
	resource = strings.TrimSpace(resource)
	if !resourceRE.MatchString(resource) {
		return fmt.Errorf("config: port_forward %q: resource %q is invalid (expected kind/name, e.g. svc/my-app)", pfName, resource)
	}
	return nil
}

func validatePortMapping(pfName, port string) error {
	port = strings.TrimSpace(port)
	if !portMappingRE.MatchString(port) {
		return fmt.Errorf("config: port_forward %q: port %q must be local:remote (e.g. 8080:80)", pfName, port)
	}

	parts := strings.Split(port, ":")
	for _, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 || n > 65535 {
			return fmt.Errorf("config: port_forward %q: port %q must use ports in range 1-65535", pfName, port)
		}
	}
	return nil
}

func validateProfileVPN(name string, profile Profile) error {
	vpn := profile.VPNType()
	switch vpn {
	case VPNTailscale:
		if strings.TrimSpace(profile.NetbirdProfile) != "" {
			return fmt.Errorf("config: profile %q: netbird_profile must not be set when vpn is tailscale", name)
		}
		if strings.TrimSpace(profile.TailscaleAccount) == "" {
			return fmt.Errorf("config: profile %q: tailscale_account is required", name)
		}
	case VPNNetbird:
		if strings.TrimSpace(profile.TailscaleAccount) != "" {
			return fmt.Errorf("config: profile %q: tailscale_account must not be set when vpn is netbird", name)
		}
		if strings.TrimSpace(profile.NetbirdProfile) == "" {
			return fmt.Errorf("config: profile %q: netbird_profile is required", name)
		}
	default:
		return fmt.Errorf("config: profile %q: vpn must be %q or %q", name, VPNTailscale, VPNNetbird)
	}
	return nil
}

func profileTargetKey(p Profile) string {
	identity := strings.TrimSpace(p.TailscaleAccount)
	if p.VPNType() == VPNNetbird {
		identity = strings.TrimSpace(p.NetbirdProfile)
	}
	return strings.TrimSpace(p.KubeContext) + "\x00" + p.VPNType() + "\x00" + identity
}

func (c *Config) UsesTailscale() bool {
	for _, p := range c.Profiles {
		if p.VPNType() == VPNTailscale {
			return true
		}
	}
	return false
}

func (c *Config) UsesNetbird() bool {
	for _, p := range c.Profiles {
		if p.VPNType() == VPNNetbird {
			return true
		}
	}
	return false
}

func (c *Config) ProfileNames() []string {
	names := make([]string, 0, len(c.Profiles))
	for name := range c.Profiles {
		names = append(names, name)
	}
	return names
}

func (c *Config) PortForwardNames() []string {
	names := make([]string, 0, len(c.PortForwards))
	for name := range c.PortForwards {
		names = append(names, name)
	}
	return names
}

func InitTemplate() string {
	return `# eggl env profiles - edit with your kubectl contexts and mesh VPN targets.
# Tailscale account refs: id from ` + "`tailscale switch --list`" + `, or tailnet slug, or account email.
# NetBird profile names: from ` + "`netbird profile list`" + `.
#
# profiles:
#   alpha:
#     kube_context: context-a
#     vpn: tailscale
#     tailscale_account: b3e1
#   beta:
#     kube_context: context-b
#     vpn: netbird
#     netbird_profile: homelab
#
# port_forwards (optional) - use with: eggl pf <name>
#   longhorn:
#     namespace: longhorn-system
#     resource: svc/longhorn-frontend
#     ports: ["8080:80"]

profiles:
  example-a:
    kube_context: context-a
    tailscale_account: account-id-a
  example-b:
    kube_context: context-b
    tailscale_account: account-id-b
`
}

func WriteInit(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config already exists: %s", path)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	if err := os.WriteFile(path, []byte(InitTemplate()), 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
