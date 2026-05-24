package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
  b:
    kube_context: ctx-b
    tailscale_account: a7f2
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Profiles) != 2 {
		t.Fatalf("profiles len = %d, want 2", len(cfg.Profiles))
	}
}

func TestValidateRequiresFields(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"bad": {KubeContext: "", TailscaleAccount: "x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoadWithPortForwards(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  a:
    kube_context: ctx-a
    tailscale_account: b3e1
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
    ports: ["8080:80"]
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	pf, ok := cfg.PortForwards["longhorn"]
	if !ok {
		t.Fatal("expected longhorn port_forward")
	}
	if pf.Namespace != "longhorn-system" || pf.Resource != "svc/longhorn-frontend" {
		t.Fatalf("port_forward = %+v", pf)
	}
}

func TestValidatePortForwardPorts(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
		PortForwards: map[string]PortForward{
			"bad": {
				Namespace: "ns",
				Resource:  "svc/x",
				Ports:     []string{"not-a-port"},
			},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for invalid port mapping")
	}
}

func TestWriteInit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "eggl", "config.yaml")
	if err := WriteInit(path); err != nil {
		t.Fatalf("WriteInit() error = %v", err)
	}
	if err := WriteInit(path); err == nil {
		t.Fatal("expected error when config exists")
	}
}

func TestDefaultPathXDG(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	got := DefaultPath()
	want := filepath.Join(dir, "eggl", "config.yaml")
	if got != want {
		t.Fatalf("DefaultPath() = %q, want %q", got, want)
	}
}

func TestDefaultPathHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("UserHomeDir:", err)
	}

	got := DefaultPath()
	want := filepath.Join(home, ".config", "eggl", "config.yaml")
	if got != want {
		t.Fatalf("DefaultPath() = %q, want %q", got, want)
	}
}

func TestProfileAndPortForwardNames(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"beta":  {KubeContext: "b", TailscaleAccount: "x"},
			"alpha": {KubeContext: "a", TailscaleAccount: "y"},
		},
		PortForwards: map[string]PortForward{
			"z": {Namespace: "ns", Resource: "svc/x"},
			"a": {Namespace: "ns", Resource: "svc/y"},
		},
	}

	profiles := cfg.ProfileNames()
	if len(profiles) != 2 {
		t.Fatalf("ProfileNames len = %d", len(profiles))
	}

	pfNames := cfg.PortForwardNames()
	if len(pfNames) != 2 {
		t.Fatalf("PortForwardNames len = %d", len(pfNames))
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(":\n\tbad"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadValidationError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("profiles: {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateDuplicateProfileTargets(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"alpha": {KubeContext: "ctx-a", TailscaleAccount: "b3e1"},
			"beta":  {KubeContext: "ctx-a", TailscaleAccount: "b3e1"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for duplicate profile targets")
	}
}

func TestValidateEmptyProfileName(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestValidatePortForwardFields(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
		PortForwards: map[string]PortForward{
			"": {Namespace: "ns", Resource: "svc/x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty port_forward name")
	}

	cfg.PortForwards = map[string]PortForward{
		"pf": {Namespace: "", Resource: "svc/x"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing namespace")
	}

	cfg.PortForwards = map[string]PortForward{
		"pf": {Namespace: "ns", Resource: ""},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing resource")
	}
}

func TestWriteInitParentNotDirectory(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "blocked")
	if err := os.WriteFile(parent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(parent, "config.yaml")

	if err := WriteInit(path); err == nil {
		t.Fatal("expected error when parent is not a directory")
	}
}

func TestLoadOversizedConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, make([]byte, MaxConfigSize+1), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for oversized config")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Fatalf("error = %v", err)
	}
}

func TestValidateInvalidNamespace(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
		PortForwards: map[string]PortForward{
			"pf": {Namespace: "INVALID", Resource: "svc/x", Ports: []string{"8080:80"}},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for invalid namespace")
	}
}

func TestValidateInvalidResource(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
		PortForwards: map[string]PortForward{
			"pf": {Namespace: "ns", Resource: "not-a-resource", Ports: []string{"8080:80"}},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for invalid resource")
	}
}

func TestValidatePortOutOfRange(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
		PortForwards: map[string]PortForward{
			"pf": {Namespace: "ns", Resource: "svc/x", Ports: []string{"0:80"}},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for port out of range")
	}
}

func TestValidateKubeContextLeadingDash(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "-bad", TailscaleAccount: "x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for kube_context starting with dash")
	}
}

func TestLoadNetbirdProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
profiles:
  homelab:
    kube_context: ctx-home
    vpn: netbird
    netbird_profile: homelab
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	p := cfg.Profiles["homelab"]
	if p.VPNType() != VPNNetbird {
		t.Fatalf("VPNType() = %q, want netbird", p.VPNType())
	}
	if !cfg.UsesNetbird() || cfg.UsesTailscale() {
		t.Fatalf("UsesNetbird=%v UsesTailscale=%v", cfg.UsesNetbird(), cfg.UsesTailscale())
	}
}

func TestValidateNetbirdRequiresProfile(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", VPN: "netbird"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for missing netbird_profile")
	}
}

func TestValidateNetbirdRejectsTailscaleAccount(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", VPN: "netbird", NetbirdProfile: "home", TailscaleAccount: "x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for tailscale_account on netbird profile")
	}
}

func TestValidateTailscaleRejectsNetbirdProfile(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x", NetbirdProfile: "home"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for netbird_profile on tailscale profile")
	}
}

func TestValidateInvalidVPN(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", VPN: "wireguard", TailscaleAccount: "x"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for invalid vpn")
	}
}

func TestValidateDuplicateNetbirdTargets(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"alpha": {KubeContext: "ctx-a", VPN: "netbird", NetbirdProfile: "home"},
			"beta":  {KubeContext: "ctx-a", VPN: "netbird", NetbirdProfile: "home"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for duplicate netbird profile targets")
	}
}

func TestUsesTailscaleDefault(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"a": {KubeContext: "ctx", TailscaleAccount: "x"},
		},
	}
	if !cfg.UsesTailscale() || cfg.UsesNetbird() {
		t.Fatalf("UsesTailscale=%v UsesNetbird=%v", cfg.UsesTailscale(), cfg.UsesNetbird())
	}
}
