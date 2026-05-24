package cmd

import (
	"fmt"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/config"
	"github.com/Robert27/eggl-cli/internal/env"
	"github.com/Robert27/eggl-cli/internal/tailscale"
	"github.com/Robert27/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var envConfigPath string

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Switch kubectl context and Tailscale account together",
	Long:  "Paired profiles in ~/.config/eggl/config.yaml. See README.md for setup.",
	Example: `  eggl env init
  eggl env show
  eggl env toggle
  eggl env use homelab`,
}

var envShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show active profile and current state",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running env show", "config", envConfigPath)

		report, err := env.Show(cmd.Context(), envOpts())
		if err != nil {
			return err
		}

		tsLabel := report.Current.TailscaleID
		if report.Current.TailscaleTailnet != "" {
			tsLabel = fmt.Sprintf("%s (%s)", report.Current.TailscaleID, report.Current.TailscaleTailnet)
		}

		profiles := make([]ui.EnvProfile, len(report.Profiles))
		for i, p := range report.Profiles {
			profiles[i] = ui.EnvProfile{
				Name:             p.Name,
				KubeContext:      p.KubeContext,
				TailscaleAccount: p.TailscaleAccount,
			}
		}

		ui.RenderEnvShow(cmd.OutOrStdout(), ui.EnvShowReport{
			ActiveProfile: report.ActiveProfile,
			Unknown:       report.Unknown,
			KubeContext:   report.Current.KubeContext,
			Tailscale:     tsLabel,
			ConfigPath:    report.ConfigPath,
			Profiles:      profiles,
		})
		return nil
	},
	SilenceUsage: true,
}

var envToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Flip to the other profile (2 profiles required)",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running env toggle", "config", envConfigPath)

		result, err := env.Toggle(cmd.Context(), envOpts())
		if err != nil {
			return err
		}

		renderEnvSwitch(cmd, result)
		return nil
	},
	SilenceUsage: true,
}

var envUseCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Apply a named profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running env use", "config", envConfigPath, "profile", args[0])

		result, err := env.Use(cmd.Context(), envOpts(), args[0])
		if err != nil {
			return err
		}

		renderEnvSwitch(cmd, result)
		return nil
	},
	SilenceUsage: true,
}

var envPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the config file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := envConfigPath
		if path == "" {
			path = config.DefaultPath()
		}
		fmt.Fprintln(cmd.OutOrStdout(), path)
		return nil
	},
	SilenceUsage: true,
}

var envInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an example config file",
	Long:  `Write a starter config to the config path if it does not exist yet.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := envConfigPath
		if path == "" {
			path = config.DefaultPath()
		}
		slog.Debug("running env init", "path", path)

		if err := config.WriteInit(path); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", path)
		return nil
	},
	SilenceUsage: true,
}

func envOpts() env.Options {
	return env.DefaultOptions(envConfigPath)
}

func renderEnvSwitch(cmd *cobra.Command, result *env.SwitchResult) {
	fromTS := result.From.TailscaleID
	if result.From.TailscaleTailnet != "" {
		fromTS = tailscale.FormatAccount(tailscale.Account{ID: result.From.TailscaleID, Tailnet: result.From.TailscaleTailnet})
	}
	toTS := result.To.TailscaleID
	if result.To.TailscaleTailnet != "" {
		toTS = tailscale.FormatAccount(tailscale.Account{ID: result.To.TailscaleID, Tailnet: result.To.TailscaleTailnet})
	}

	ui.RenderEnvSwitch(cmd.OutOrStdout(), ui.EnvSwitchResult{
		FromProfile: result.FromProfile,
		ToProfile:   result.ToProfile,
		FromKube:    result.From.KubeContext,
		ToKube:      result.To.KubeContext,
		FromTS:      fromTS,
		ToTS:        toTS,
	})
}

func init() {
	envCmd.PersistentFlags().StringVar(&envConfigPath, "config", "", "Path to config file")

	envCmd.AddCommand(envShowCmd, envToggleCmd, envUseCmd, envPathCmd, envInitCmd)
	rootCmd.AddCommand(envCmd)
}
