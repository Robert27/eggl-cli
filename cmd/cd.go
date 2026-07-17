package cmd

import (
	"fmt"
	"log/slog"

	"github.com/roberteggl/eggl-cli/internal/cd"
	"github.com/roberteggl/eggl-cli/internal/config"
	"github.com/spf13/cobra"
)

var cdConfigPath string

var cdCmd = &cobra.Command{
	Use:   "cd <alias>",
	Short: "Print a configured directory path",
	Long:  "Named directories in ~/.config/eggl/config.yaml. Use: cd \"$(eggl cd <alias>)\"",
	Example: `  cd "$(eggl cd homelab)"
  eggl cd list`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		slog.Debug("running cd", "alias", args[0], "config", cdConfigPath)

		path, err := cd.Resolve(args[0], cd.DefaultOptions(cdConfigPath))
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), path)
		return nil
	},
	SilenceUsage: true,
}

var cdListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running cd list", "config", cdConfigPath)

		path := cdConfigPath
		if path == "" {
			path = config.DefaultPath()
		}
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		entries := cd.List(cfg)
		if len(entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no directories configured")
			return nil
		}

		for _, e := range entries {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", e.Name, e.Path)
		}
		return nil
	},
	SilenceUsage: true,
}

func cdComplete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	path := cdConfigPath
	if path == "" {
		path = config.DefaultPath()
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return cfg.DirectoryNames(), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	cdCmd.PersistentFlags().StringVar(&cdConfigPath, "config", "", "Path to config file")
	cdCmd.ValidArgsFunction = cdComplete

	cdCmd.AddCommand(cdListCmd)
	rootCmd.AddCommand(cdCmd)
}
