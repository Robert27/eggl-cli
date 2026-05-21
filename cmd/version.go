package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
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
	versionCmd.Flags().BoolVar(&versionShort, "short", false, "Print only the version number")
	rootCmd.AddCommand(versionCmd)
}
