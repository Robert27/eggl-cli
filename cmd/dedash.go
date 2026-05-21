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
	Long:  `Recursively scan a directory and replace Unicode em-dashes (-) with ASCII hyphens (-) in text files, skipping binaries and common non-text paths.`,
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
	dedashCmd.Flags().StringVar(&dedashPath, "path", ".", "Root directory to scan")
	dedashCmd.Flags().BoolVar(&dedashDryRun, "dry-run", false, "Preview changes without writing files")
	rootCmd.AddCommand(dedashCmd)
}
