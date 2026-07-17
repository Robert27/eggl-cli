package cmd

import (
	"fmt"
	"log/slog"

	"github.com/roberteggl/eggl-cli/internal/eol"
	"github.com/roberteggl/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	eolPath          string
	eolDryRun        bool
	eolYes           bool
	eolExtensions    []string
	eolIncludeHidden bool
	eolGitDiff       bool
	eolDiffBase      string
)

var eolCmd = &cobra.Command{
	Use:   "eol",
	Short: "Normalize line endings to LF in text files",
	Long: `Recursively scan a directory and convert CRLF or CR line endings to LF in text files.

Skips binaries, dotfiles and dot-directories (unless --include-hidden), and common non-text directories such as node_modules and .git.
Files larger than 50 MiB are skipped.
Use --ext to limit processing to specific file extensions (for example md,txt or .md,.txt).
Use --diff to limit processing to files with staged or unstaged git changes (must be inside a git repo).
Use --diff-base to limit processing to files changed on the current branch since a ref (e.g. main). Deletions are skipped.
Prints a one-line summary plus a list of changed files. Use --dry-run to preview without writing.

When not using --dry-run, you are prompted to confirm before any files are modified.
Use --yes to skip the prompt (required in non-interactive environments).

With --verbose, stderr logs the scan root, skipped paths, and each file that would be or was modified.`,
	Example: `  eggl eol --dry-run
  eggl eol --path ./docs --yes
  eggl eol --ext md,txt --dry-run
  eggl eol --diff --dry-run
  eggl eol --diff-base main --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if eolGitDiff && eolDiffBase != "" {
			return fmt.Errorf("cannot use --diff with --diff-base")
		}

		slog.Debug("running eol", "path", eolPath, "dry_run", eolDryRun, "yes", eolYes, "extensions", eolExtensions, "git_diff", eolGitDiff, "diff_base", eolDiffBase)

		report, err := eol.Run(cmd.Context(), eol.Options{
			Root:          eolPath,
			DryRun:        eolDryRun,
			Yes:           eolYes,
			Extensions:    eolExtensions,
			IncludeHidden: eolIncludeHidden,
			GitDiff:       eolGitDiff,
			GitDiffBase:   eolDiffBase,
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

		changes := make([]ui.EOLChange, len(report.Changes))
		for i, change := range report.Changes {
			changes[i] = ui.EOLChange{
				Path:         change.Path,
				Replacements: change.Replacements,
			}
		}
		ui.RenderEOL(cmd.OutOrStdout(), ui.EOLSummary{
			Scanned:           report.Scanned,
			Modified:          report.Modified,
			Skipped:           report.Skipped,
			TotalReplacements: eol.TotalReplacements(report),
			Changes:           changes,
			DryRun:            eolDryRun,
		})
		return nil
	},
	SilenceUsage: true,
}

func init() {
	eolCmd.Flags().StringVar(&eolPath, "path", ".", "Directory tree to scan (default: current directory)")
	eolCmd.Flags().BoolVar(&eolDryRun, "dry-run", false, "Report changes without writing files")
	eolCmd.Flags().BoolVarP(&eolYes, "yes", "y", false, "Confirm writes without prompting (required when stdin is not a terminal)")
	eolCmd.Flags().StringSliceVar(&eolExtensions, "ext", nil, "Only process files with these extensions (comma-separated or repeatable, e.g. md,txt)")
	eolCmd.Flags().BoolVar(&eolIncludeHidden, "include-hidden", false, "Process dotfiles and dot-directories")
	eolCmd.Flags().BoolVar(&eolGitDiff, "diff", false, "Only process files with staged or unstaged git changes (skips deletions)")
	eolCmd.Flags().StringVar(&eolDiffBase, "diff-base", "", "Only process files changed on the current branch since a git ref (e.g. main; skips deletions)")
	rootCmd.AddCommand(eolCmd)
}
