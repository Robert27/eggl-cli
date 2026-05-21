package cmd

import (
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/dedash"
	"github.com/Robert27/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	dedashPath   string
	dedashDryRun bool
)

var dedashCmd = &cobra.Command{
	Use:   "dedash",
	Short: "Replace em-dashes with hyphens in text files",
	Long: `Recursively scan a directory and replace Unicode em-dashes (—) with ASCII hyphens (-) in text files.

Skips binaries, hidden paths, and common non-text directories such as node_modules and .git.
Prints a one-line summary plus a list of changed files. Use --dry-run to preview without writing.

With --verbose, stderr logs the scan root, skipped paths, and each file that would be or was modified.`,
	Example: `  eggl dedash --dry-run
  eggl dedash --path ./docs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running dedash", "path", dedashPath, "dry_run", dedashDryRun)

		report, err := dedash.Run(cmd.Context(), dedash.Options{
			Root:   dedashPath,
			DryRun: dedashDryRun,
		})
		if err != nil {
			return err
		}

		changes := make([]ui.DedashChange, len(report.Changes))
		for i, change := range report.Changes {
			changes[i] = ui.DedashChange{
				Path:         change.Path,
				Replacements: change.Replacements,
			}
		}
		ui.RenderDedash(cmd.OutOrStdout(), ui.DedashSummary{
			Scanned:           report.Scanned,
			Modified:          report.Modified,
			Skipped:           report.Skipped,
			TotalReplacements: dedash.TotalReplacements(report),
			Changes:           changes,
			DryRun:            dedashDryRun,
		})
		return nil
	},
	SilenceUsage: true,
}

func init() {
	dedashCmd.Flags().StringVar(&dedashPath, "path", ".", "Directory tree to scan (default: current directory)")
	dedashCmd.Flags().BoolVar(&dedashDryRun, "dry-run", false, "Report changes without writing files")
	rootCmd.AddCommand(dedashCmd)
}
