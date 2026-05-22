package cmd

import (
	"fmt"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/git"
	"github.com/spf13/cobra"
)

var emptyPush bool

var emptyCmd = &cobra.Command{
	Use:   "empty",
	Short: "Create an empty git commit",
	Long: `Create an empty git commit in the current repository.

Useful for retriggering CI pipelines without changing files.
Use --push to push the commit after creating it.`,
	Example: `  eggl empty
  eggl empty --push
  eggl empty -p`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running empty commit", "push", emptyPush)

		result, err := git.Run(cmd.Context(), git.Options{
			Push: emptyPush,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "empty commit %s\n", result.Hash)
		if result.Pushed {
			fmt.Fprintln(cmd.OutOrStdout(), "pushed")
		}
		return nil
	},
	SilenceUsage: true,
}

func init() {
	emptyCmd.Flags().BoolVarP(&emptyPush, "push", "p", false, "Push the commit after creating it")
	rootCmd.AddCommand(emptyCmd)
}
