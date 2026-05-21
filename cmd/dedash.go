package cmd

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/dedash"
	"github.com/spf13/cobra"
)

var (
	dedashPath   string
	dedashDryRun bool
)

var dedashCmd = &cobra.Command{
	Use:   "dedash",
	Short: "Replace em-dashes with hyphens in text files",
	Long:  `Recursively scan a directory and replace Unicode em-dashes (—) with ASCII hyphens (-) in text files, skipping binaries and common non-text paths.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running dedash", "path", dedashPath, "dry_run", dedashDryRun)

		report, err := dedash.Run(cmd.Context(), dedash.Options{
			Root:   dedashPath,
			DryRun: dedashDryRun,
		})
		if err != nil {
			return err
		}

		renderDedashReport(cmd.OutOrStdout(), report, dedashDryRun)
		return nil
	},
	SilenceUsage: true,
}

func init() {
	dedashCmd.Flags().StringVar(&dedashPath, "path", ".", "Root directory to scan")
	dedashCmd.Flags().BoolVar(&dedashDryRun, "dry-run", false, "Preview changes without writing files")
	rootCmd.AddCommand(dedashCmd)
}

func renderDedashReport(w io.Writer, report *dedash.Report, dryRun bool) {
	total := dedash.TotalReplacements(report)

	if dryRun {
		fmt.Fprintf(w, "dry-run: scanned %d files, would modify %d (%d replacements)\n",
			report.Scanned, report.Modified, total)
	} else {
		fmt.Fprintf(w, "Scanned %d files, modified %d (%d replacements)\n",
			report.Scanned, report.Modified, total)
	}

	for _, change := range report.Changes {
		fmt.Fprintf(w, "  %s (%d)\n", change.Path, change.Replacements)
	}
}
