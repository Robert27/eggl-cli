package cmd

import (
	"fmt"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionShort bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long: `Show the installed eggl version, git commit, and build date.

Use --short for script-friendly output containing only the version string.`,
	Example: `  eggl version
  eggl version --short`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("printing version", "short", versionShort)

		if versionShort {
			fmt.Fprintln(cmd.OutOrStdout(), version)
			return nil
		}

		ui.RenderVersion(cmd.OutOrStdout(), ui.VersionInfo{
			Version: version,
			Commit:  commit,
			Date:    date,
		})
		return nil
	},
	SilenceUsage: true,
}

func init() {
	versionCmd.Flags().BoolVar(&versionShort, "short", false, "Print only the version string (no commit or build date)")
	rootCmd.AddCommand(versionCmd)
}
