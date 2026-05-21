package cmd

import (
	"fmt"

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

		fmt.Fprintf(cmd.OutOrStdout(), "eggl version %s\n", version)
		fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", commit)
		fmt.Fprintf(cmd.OutOrStdout(), "built:  %s\n", date)
		return nil
	},
	SilenceUsage: true,
}

func init() {
	versionCmd.Flags().BoolVar(&versionShort, "short", false, "Print only the version number")
	rootCmd.AddCommand(versionCmd)
}
