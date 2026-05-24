package env

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/Robert27/eggl-cli/internal/config"
	"github.com/Robert27/eggl-cli/internal/kube"
	"github.com/Robert27/eggl-cli/internal/netbird"
	"github.com/Robert27/eggl-cli/internal/tailscale"
)

type Options struct {
	ConfigPath string
	Kube       kube.Runner
	TS         tailscale.Runner
	NB         netbird.Runner
}

type State struct {
	KubeContext      string
	TailscaleID      string
	TailscaleTailnet string
	NetbirdProfile   string
}

type meshRead struct {
	accounts []tailscale.Account
	profiles []netbird.Profile
}

type Report struct {
	ActiveProfile string
	Unknown       bool
	Current       State
	ConfigPath    string
	Profiles      []ProfileInfo
	ShowTailscale bool
	ShowNetbird   bool
}

type ProfileInfo struct {
	Name        string
	KubeContext string
	Mesh        string
}

type SwitchResult struct {
	FromProfile string
	ToProfile   string
	From        State
	To          State
	MeshVPN     string
}

func DefaultOptions(configPath string) Options {
	return Options{
		ConfigPath: configPath,
		Kube:       kube.CLI{},
		TS:         tailscale.CLI{},
		NB:         netbird.CLI{},
	}
}

func Show(ctx context.Context, opts Options) (*Report, error) {
	cfg, path, err := loadConfig(opts.ConfigPath)
	if err != nil {
		return nil, err
	}

	state, mr, err := readState(ctx, opts, cfg)
	if err != nil {
		return nil, err
	}

	report := &Report{
		Current:       state,
		ConfigPath:    path,
		Profiles:      profileInfos(cfg),
		ShowTailscale: cfg.UsesTailscale(),
		ShowNetbird:   cfg.UsesNetbird(),
	}

	active, unknown, err := detectProfile(cfg, state, mr)
	if err != nil {
		return nil, err
	}
	report.ActiveProfile = active
	report.Unknown = unknown
	return report, nil
}

func Use(ctx context.Context, opts Options, profileName string) (*SwitchResult, error) {
	cfg, _, err := loadConfig(opts.ConfigPath)
	if err != nil {
		return nil, err
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("profile %q not found in config", profileName)
	}

	return applyProfile(ctx, opts, cfg, profileName, profile)
}

func Toggle(ctx context.Context, opts Options) (*SwitchResult, error) {
	cfg, _, err := loadConfig(opts.ConfigPath)
	if err != nil {
		return nil, err
	}

	if len(cfg.Profiles) != 2 {
		return nil, fmt.Errorf("toggle requires exactly 2 profiles in config (found %d)", len(cfg.Profiles))
	}

	state, mr, err := readState(ctx, opts, cfg)
	if err != nil {
		return nil, err
	}

	active, unknown, err := detectProfile(cfg, state, mr)
	if err != nil {
		return nil, err
	}
	if unknown {
		return nil, fmt.Errorf("current kube context and mesh VPN state do not match any profile; use `eggl env use <name>` (kube=%q, %s)",
			state.KubeContext, formatMeshState(cfg, state, mr))
	}

	names := sortedProfileNames(cfg)
	var target string
	for _, name := range names {
		if name != active {
			target = name
			break
		}
	}

	return applyProfile(ctx, opts, cfg, target, cfg.Profiles[target])
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

func readState(ctx context.Context, opts Options, cfg *config.Config) (State, meshRead, error) {
	kubeCtx, err := opts.Kube.CurrentContext(ctx)
	if err != nil {
		return State{}, meshRead{}, err
	}

	state := State{KubeContext: kubeCtx}
	var mr meshRead

	if cfg.UsesTailscale() {
		accounts, err := opts.TS.ListAccounts(ctx)
		if err != nil {
			return State{}, meshRead{}, err
		}
		current, err := tailscale.CurrentAccount(accounts)
		if err != nil {
			return State{}, meshRead{}, err
		}
		state.TailscaleID = current.ID
		state.TailscaleTailnet = current.Tailnet
		mr.accounts = accounts
	}

	if cfg.UsesNetbird() {
		profiles, err := opts.NB.ListProfiles(ctx)
		if err != nil {
			return State{}, meshRead{}, err
		}
		current, err := netbird.CurrentProfile(profiles)
		if err != nil {
			return State{}, meshRead{}, err
		}
		state.NetbirdProfile = current.Name
		mr.profiles = profiles
	}

	return state, mr, nil
}

func detectProfile(cfg *config.Config, state State, mr meshRead) (string, bool, error) {
	var matches []string
	for name, profile := range cfg.Profiles {
		ok, err := profileMatches(profile, state, mr)
		if err != nil {
			continue
		}
		if ok {
			matches = append(matches, name)
		}
	}

	switch len(matches) {
	case 0:
		return "", true, nil
	case 1:
		return matches[0], false, nil
	default:
		sort.Strings(matches)
		return "", false, fmt.Errorf("ambiguous profile match for current state: %s", strings.Join(matches, ", "))
	}
}

func profileMatches(profile config.Profile, state State, mr meshRead) (bool, error) {
	if profile.KubeContext != state.KubeContext {
		return false, nil
	}

	switch profile.VPNType() {
	case config.VPNTailscale:
		tsAccount, err := tailscale.ResolveAccount(profile.TailscaleAccount, mr.accounts)
		if err != nil {
			return false, err
		}
		return tsAccount.ID == state.TailscaleID, nil
	case config.VPNNetbird:
		nbProfile, err := netbird.ResolveProfile(profile.NetbirdProfile, mr.profiles)
		if err != nil {
			return false, err
		}
		return strings.EqualFold(nbProfile.Name, state.NetbirdProfile), nil
	default:
		return false, fmt.Errorf("unsupported vpn %q", profile.VPNType())
	}
}

func applyProfile(ctx context.Context, opts Options, cfg *config.Config, name string, profile config.Profile) (*SwitchResult, error) {
	state, mr, err := readState(ctx, opts, cfg)
	if err != nil {
		return nil, err
	}

	active, _, err := detectProfile(cfg, state, mr)
	if err != nil {
		return nil, err
	}

	result := &SwitchResult{
		FromProfile: active,
		ToProfile:   name,
		From:        state,
		To: State{
			KubeContext: profile.KubeContext,
		},
		MeshVPN: profile.VPNType(),
	}

	var applied []string

	if state.KubeContext != profile.KubeContext {
		slog.Debug("switching kube context", "from", state.KubeContext, "to", profile.KubeContext)
		if err := opts.Kube.UseContext(ctx, profile.KubeContext); err != nil {
			return result, err
		}
		applied = append(applied, "kube")
	}

	switch profile.VPNType() {
	case config.VPNTailscale:
		targetTS, err := tailscale.ResolveAccount(profile.TailscaleAccount, mr.accounts)
		if err != nil {
			return result, err
		}
		result.To.TailscaleID = targetTS.ID
		result.To.TailscaleTailnet = targetTS.Tailnet
		if state.TailscaleID != targetTS.ID {
			slog.Debug("switching tailscale account", "from", state.TailscaleID, "to", targetTS.ID)
			if err := opts.TS.Switch(ctx, targetTS.ID); err != nil {
				if len(applied) > 0 {
					return result, fmt.Errorf("%w (kube context already switched to %q)", err, profile.KubeContext)
				}
				return result, err
			}
			applied = append(applied, "tailscale")
		}
	case config.VPNNetbird:
		targetNB, err := netbird.ResolveProfile(profile.NetbirdProfile, mr.profiles)
		if err != nil {
			return result, err
		}
		result.To.NetbirdProfile = targetNB.Name
		if !strings.EqualFold(state.NetbirdProfile, targetNB.Name) {
			slog.Debug("switching netbird profile", "from", state.NetbirdProfile, "to", targetNB.Name)
			if err := opts.NB.SelectProfile(ctx, targetNB.Name); err != nil {
				if len(applied) > 0 {
					return result, fmt.Errorf("%w (kube context already switched to %q)", err, profile.KubeContext)
				}
				return result, err
			}
			applied = append(applied, "netbird")
		}
	default:
		return result, fmt.Errorf("unsupported vpn %q in profile %q", profile.VPNType(), name)
	}

	if len(applied) == 0 {
		slog.Debug("profile already active", "profile", name)
	}

	return result, nil
}

func profileInfos(cfg *config.Config) []ProfileInfo {
	names := sortedProfileNames(cfg)
	out := make([]ProfileInfo, 0, len(names))
	for _, name := range names {
		p := cfg.Profiles[name]
		out = append(out, ProfileInfo{
			Name:        name,
			KubeContext: p.KubeContext,
			Mesh:        profileMeshLabel(p),
		})
	}
	return out
}

func profileMeshLabel(p config.Profile) string {
	switch p.VPNType() {
	case config.VPNNetbird:
		return "netbird:" + p.NetbirdProfile
	default:
		return "tailscale:" + p.TailscaleAccount
	}
}

func formatMeshState(cfg *config.Config, state State, mr meshRead) string {
	var parts []string
	if cfg.UsesTailscale() {
		parts = append(parts, "tailscale="+tailscale.FormatAccount(findCurrentAccount(mr.accounts)))
	}
	if cfg.UsesNetbird() && state.NetbirdProfile != "" {
		parts = append(parts, "netbird="+state.NetbirdProfile)
	}
	return strings.Join(parts, ", ")
}

func sortedProfileNames(cfg *config.Config) []string {
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func findCurrentAccount(accounts []tailscale.Account) tailscale.Account {
	a, err := tailscale.CurrentAccount(accounts)
	if err != nil {
		return tailscale.Account{}
	}
	return a
}
