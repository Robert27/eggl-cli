package cmd

import (
	"fmt"
	"log/slog"

	"github.com/roberteggl/eggl-cli/internal/dedash"
	"github.com/roberteggl/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	dedashPath          string
	dedashDryRun        bool
	dedashYes           bool
	dedashExtensions    []string
	dedashIncludeHidden bool
	dedashGitDiff       bool
	dedashDiffBase      string
)

var dedashCmd = &cobra.Command{
	Use:   "dedash",
	Short: "Replace em-dashes with hyphens in text files",
	Long: `Recursively scan a directory and replace Unicode em-dashes (—) with ASCII hyphens (-) in text files.

Skips binaries, dotfiles and dot-directories (unless --include-hidden), and common non-text directories such as node_modules and .git.
Files larger than 50 MiB are skipped.
Use --ext to limit processing to specific file extensions (for example md,txt or .md,.txt).
Use --diff to limit processing to files with staged or unstaged git changes (must be inside a git repo).
Use --diff-base to limit processing to files changed on the current branch since a ref (e.g. main). Deletions are skipped.
Prints a one-line summary plus a list of changed files. Use --dry-run to preview without writing.

When not using --dry-run, you are prompted to confirm before any files are modified.
Use --yes to skip the prompt (required in non-interactive environments).

With --verbose, stderr logs the scan root, skipped paths, and each file that would be or was modified.`,
	Example: `  eggl dedash --dry-run
  eggl dedash --path ./docs --yes
  eggl dedash --ext md,txt --dry-run
  eggl dedash --diff --dry-run
  eggl dedash --diff-base main --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dedashGitDiff && dedashDiffBase != "" {
			return fmt.Errorf("cannot use --diff with --diff-base")
		}

		slog.Debug("running dedash", "path", dedashPath, "dry_run", dedashDryRun, "yes", dedashYes, "extensions", dedashExtensions, "git_diff", dedashGitDiff, "diff_base", dedashDiffBase)

		report, err := dedash.Run(cmd.Context(), dedash.Options{
			Root:          dedashPath,
			DryRun:        dedashDryRun,
			Yes:           dedashYes,
			Extensions:    dedashExtensions,
			IncludeHidden: dedashIncludeHidden,
			GitDiff:       dedashGitDiff,
			GitDiffBase:   dedashDiffBase,
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
	dedashCmd.Flags().BoolVar(&dedashGitDiff, "diff", false, "Only process files with staged or unstaged git changes (skips deletions)")
	dedashCmd.Flags().StringVar(&dedashDiffBase, "diff-base", "", "Only process files changed on the current branch since a git ref (e.g. main; skips deletions)")
	rootCmd.AddCommand(dedashCmd)
}
