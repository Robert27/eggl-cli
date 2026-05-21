package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "eggl",
	Short: "A general-purpose helper CLI",
	Long:  `eggl is a personal helper CLI for dev workflow, cloud, and everyday tasks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initLogging(verbose, cmd.ErrOrStderr())
		return nil
	},
}

func initLogging(verbose bool, w io.Writer) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable debug logging")
}
