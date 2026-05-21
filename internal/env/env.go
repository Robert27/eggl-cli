package env

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/Robert27/eggl-cli/internal/config"
	"github.com/Robert27/eggl-cli/internal/kube"
	"github.com/Robert27/eggl-cli/internal/tailscale"
)

type Options struct {
	ConfigPath string
	Kube       kube.Runner
	TS         tailscale.Runner
}

type State struct {
	KubeContext      string
	TailscaleID      string
	TailscaleTailnet string
}

type Report struct {
	ActiveProfile string
	Unknown       bool
	Current       State
	ConfigPath    string
	Profiles      []ProfileInfo
}

type ProfileInfo struct {
	Name             string
	KubeContext      string
	TailscaleAccount string
}

type SwitchResult struct {
	FromProfile string
	ToProfile   string
	From        State
	To          State
}

func DefaultOptions(configPath string) Options {
	return Options{
		ConfigPath: configPath,
		Kube:       kube.CLI{},
		TS:         tailscale.CLI{},
	}
}

func Show(ctx context.Context, opts Options) (*Report, error) {
	cfg, path, err := loadConfig(opts.ConfigPath)
	if err != nil {
		return nil, err
	}

	state, accounts, err := readState(ctx, opts)
	if err != nil {
		return nil, err
	}

	report := &Report{
		Current:    state,
		ConfigPath: path,
		Profiles:   profileInfos(cfg),
	}

	report.ActiveProfile, report.Unknown = detectProfile(cfg, state, accounts)
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

	state, accounts, err := readState(ctx, opts)
	if err != nil {
		return nil, err
	}

	active, unknown := detectProfile(cfg, state, accounts)
	if unknown {
		return nil, fmt.Errorf("current kube context and tailscale account do not match any profile; use `eggl env use <name>` (kube=%q, tailscale=%s)",
			state.KubeContext, tailscale.FormatAccount(findCurrentAccount(accounts)))
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

func readState(ctx context.Context, opts Options) (State, []tailscale.Account, error) {
	kubeCtx, err := opts.Kube.CurrentContext(ctx)
	if err != nil {
		return State{}, nil, err
	}

	accounts, err := opts.TS.ListAccounts(ctx)
	if err != nil {
		return State{}, nil, err
	}

	current, err := tailscale.CurrentAccount(accounts)
	if err != nil {
		return State{}, accounts, err
	}

	return State{
		KubeContext:      kubeCtx,
		TailscaleID:      current.ID,
		TailscaleTailnet: current.Tailnet,
	}, accounts, nil
}

func detectProfile(cfg *config.Config, state State, accounts []tailscale.Account) (string, bool) {
	for name, profile := range cfg.Profiles {
		tsAccount, err := tailscale.ResolveAccount(profile.TailscaleAccount, accounts)
		if err != nil {
			continue
		}
		if profile.KubeContext == state.KubeContext && tsAccount.ID == state.TailscaleID {
			return name, false
		}
	}
	return "", true
}

func applyProfile(ctx context.Context, opts Options, cfg *config.Config, name string, profile config.Profile) (*SwitchResult, error) {
	state, accounts, err := readState(ctx, opts)
	if err != nil {
		return nil, err
	}

	active, _ := detectProfile(cfg, state, accounts)

	targetTS, err := tailscale.ResolveAccount(profile.TailscaleAccount, accounts)
	if err != nil {
		return nil, err
	}

	result := &SwitchResult{
		FromProfile: active,
		ToProfile:   name,
		From:        state,
		To: State{
			KubeContext:      profile.KubeContext,
			TailscaleID:      targetTS.ID,
			TailscaleTailnet: targetTS.Tailnet,
		},
	}

	var applied []string

	if state.TailscaleID != targetTS.ID {
		slog.Debug("switching tailscale account", "from", state.TailscaleID, "to", targetTS.ID)
		if err := opts.TS.Switch(ctx, targetTS.ID); err != nil {
			return result, err
		}
		applied = append(applied, "tailscale")
	}

	if state.KubeContext != profile.KubeContext {
		slog.Debug("switching kube context", "from", state.KubeContext, "to", profile.KubeContext)
		if err := opts.Kube.UseContext(ctx, profile.KubeContext); err != nil {
			if len(applied) > 0 {
				return result, fmt.Errorf("%w (tailscale already switched to %s)", err, targetTS.ID)
			}
			return result, err
		}
		applied = append(applied, "kube")
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
			Name:             name,
			KubeContext:      p.KubeContext,
			TailscaleAccount: p.TailscaleAccount,
		})
	}
	return out
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
