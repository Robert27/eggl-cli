package cmd

import (
	"errors"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/doctor"
	"github.com/Robert27/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCheckPath string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check local environment and dependencies",
	Long: `Run quick sanity checks for tools and paths eggl relies on.

Checks the Go runtime, platform (GOOS/GOARCH), home directory, kubectl, git,
tailscale on PATH, and eggl config validity (when present).
Use --check-path to validate a different directory instead of $HOME.
Exits with an error when any check fails.

With --verbose, stderr logs each check name and result as it runs.`,
	Example: `  eggl doctor
  eggl doctor --check-path /tmp`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running doctor checks", "check_path", doctorCheckPath)

		report, err := doctor.Run(cmd.Context(), doctor.Options{
			CheckPath: doctorCheckPath,
		})
		if err != nil {
			return err
		}

		checks := make([]ui.DoctorCheck, len(report.Checks))
		for i, check := range report.Checks {
			checks[i] = ui.DoctorCheck{
				Name:   check.Name,
				Status: check.Status,
				Detail: check.Detail,
				OK:     check.OK,
			}
		}

		ui.RenderDoctor(cmd.OutOrStdout(), checks)

		if doctor.HasFailures(report) {
			return errors.New("one or more checks failed")
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	doctorCmd.Flags().StringVar(&doctorCheckPath, "check-path", "", "Directory to validate instead of $HOME (must exist and be readable)")
	rootCmd.AddCommand(doctorCmd)
}
