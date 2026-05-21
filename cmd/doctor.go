package cmd

import (
	"errors"
	"log/slog"

	"github.com/Robert27/eggl-cli/internal/doctor"
	"github.com/spf13/cobra"
)

var doctorCheckPath string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check local environment and dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("running doctor checks", "check_path", doctorCheckPath)

		report, err := doctor.Run(cmd.Context(), doctor.Options{
			CheckPath: doctorCheckPath,
		})
		if err != nil {
			return err
		}

		doctor.Print(cmd.OutOrStdout(), report)

		if doctor.HasFailures(report) {
			return errors.New("one or more checks failed")
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	doctorCmd.Flags().StringVar(&doctorCheckPath, "check-path", "", "Path to verify instead of $HOME")
	rootCmd.AddCommand(doctorCmd)
}
