package cmd

import (
	"fmt"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/dedash"
	"github.com/Robert27/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	dedashPath          string
	dedashDryRun        bool
	dedashYes           bool
	dedashExtensions    []string
	dedashIncludeHidden bool
)

var dedashCmd = &cobra.Command{
	Use:   "dedash",
	Short: "Replace em-dashes with hyphens in text files",
	Long: `Recursively scan a directory and replace Unicode em-dashes (—) with ASCII hyphens (-) in text files.

Skips binaries, dotfiles and dot-directories (unless --include-hidden), and common non-text directories such as node_modules and .git.
Files larger than 50 MiB are skipped.
Use --ext to limit processing to specific file extensions (for example md,txt or .md,.txt).
Prints a one-line summary plus a list of changed files. Use --dry-run to preview without writing.

When not using --dry-run, you are prompted to confirm before any files are modified.
Use --yes to skip the prompt (required in non-interactive environments).

With --verbose, stderr logs the scan root, skipped paths, and each file that would be or was modified.`,
	Example: `  eggl dedash --dry-run
  eggl dedash --path ./docs --yes
  eggl dedash --ext md,txt --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running dedash", "path", dedashPath, "dry_run", dedashDryRun, "yes", dedashYes, "extensions", dedashExtensions)

		report, err := dedash.Run(cmd.Context(), dedash.Options{
			Root:          dedashPath,
			DryRun:        dedashDryRun,
			Yes:           dedashYes,
			Extensions:    dedashExtensions,
			IncludeHidden: dedashIncludeHidden,
			Input:         cmd.InOrStdin(),
			Output:        cmd.ErrOrStderr(),
		})
		if err != nil {
			return err
		}

		if report.Cancelled {
			fmt.Fprintln(cmd.ErrOrStderr(), "cancelled")
			return nil
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
	dedashCmd.Flags().BoolVarP(&dedashYes, "yes", "y", false, "Confirm writes without prompting (required when stdin is not a terminal)")
	dedashCmd.Flags().StringSliceVar(&dedashExtensions, "ext", nil, "Only process files with these extensions (comma-separated or repeatable, e.g. md,txt)")
	dedashCmd.Flags().BoolVar(&dedashIncludeHidden, "include-hidden", false, "Process dotfiles and dot-directories")
	rootCmd.AddCommand(dedashCmd)
}
