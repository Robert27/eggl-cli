package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Robert27/eggl-cli/internal/config"
	"github.com/Robert27/eggl-cli/internal/pf"
	"github.com/spf13/cobra"
)

var pfConfigPath string

var pfCmd = &cobra.Command{
	Use:   "pf [service]",
	Short: "Port-forward configured Kubernetes services",
	Long:  "Named port-forwards in ~/.config/eggl/config.yaml. Uses the active kubectl context.",
	Example: `  eggl pf list
  eggl pf longhorn`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		slog.Debug("running pf", "service", args[0], "config", pfConfigPath)
		return pf.Run(cmd.Context(), args[0], pf.DefaultOptions(pfConfigPath))
	},
	SilenceUsage: true,
}

var pfListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured port-forwards",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running pf list", "config", pfConfigPath)

		path := pfConfigPath
		if path == "" {
			path = config.DefaultPath()
		}
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}

		entries := pf.List(cfg)
		if len(entries) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no port_forwards configured")
			return nil
		}

		for _, e := range entries {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s/%s\t%s\n",
				e.Name, e.Namespace, e.Resource, strings.Join(e.Ports, ","))
		}
		return nil
	},
	SilenceUsage: true,
}

func pfComplete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	path := pfConfigPath
	if path == "" {
		path = config.DefaultPath()
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := cfg.PortForwardNames()
	return names, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	pfCmd.PersistentFlags().StringVar(&pfConfigPath, "config", "", "Path to config file")
	pfCmd.ValidArgsFunction = pfComplete

	pfCmd.AddCommand(pfListCmd)
	rootCmd.AddCommand(pfCmd)
}
