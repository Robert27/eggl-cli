package cmd

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/roberteggl/eggl-cli/internal/kill"
	"github.com/spf13/cobra"
)

var (
	killDryRun bool
	killYes    bool
	killForce  bool
)

var killCmd = &cobra.Command{
	Use:   "kill <port>",
	Short: "Kill processes listening on a TCP port",
	Long: `Find and terminate processes that are listening on the given local TCP port.

Useful when a stale port-forward or dev server still holds the port.
Use --dry-run to list matching processes without sending a signal.
Use --yes to skip the confirmation prompt (required in non-interactive environments).
Use --force to send SIGKILL instead of SIGTERM.

With --verbose, stderr logs the target port and each process that would be or was killed.`,
	Example: `  eggl kill 8080
  eggl kill --dry-run 3000
  eggl kill --yes --force 8080`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid port %q", args[0])
		}

		slog.Debug("running kill", "port", port, "dry_run", killDryRun, "yes", killYes, "force", killForce)

		result, err := kill.Run(cmd.Context(), kill.Options{
			Port:   port,
			DryRun: killDryRun,
			Yes:    killYes,
			Force:  killForce,
			Input:  cmd.InOrStdin(),
			Output: cmd.ErrOrStderr(),
		})
		if err != nil {
			return err
		}

		if result.Cancelled {
			fmt.Fprintln(cmd.ErrOrStderr(), "cancelled")
			return nil
		}

		processes := result.Found
		if killDryRun {
			processes = result.Found
		} else if len(result.Killed) > 0 {
			processes = result.Killed
		}

		if len(processes) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "no process listening on port %d\n", port)
			return nil
		}

		for _, proc := range processes {
			line := fmt.Sprintf("pid %d", proc.PID)
			if proc.Name != "" {
				line = fmt.Sprintf("%s (%s)", line, proc.Name)
			}
			if killDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "would kill %s\n", line)
				continue
			}
			if killForce {
				fmt.Fprintf(cmd.OutOrStdout(), "killed %s (SIGKILL)\n", line)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "killed %s\n", line)
			}
		}
		return nil
	},
	SilenceUsage: true,
}

func init() {
	killCmd.Flags().BoolVar(&killDryRun, "dry-run", false, "List matching processes without killing them")
	killCmd.Flags().BoolVarP(&killYes, "yes", "y", false, "Confirm kills without prompting (required when stdin is not a terminal)")
	killCmd.Flags().BoolVarP(&killForce, "force", "f", false, "Send SIGKILL instead of SIGTERM")
	rootCmd.AddCommand(killCmd)
}
