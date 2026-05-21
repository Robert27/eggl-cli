package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/Robert27/eggl-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "eggl",
	Short: "A general-purpose helper CLI",
	Long: `eggl is a personal helper CLI for dev workflow, cloud, and everyday tasks.

Use --verbose to print operation details (skipped paths, check progress, and similar) to stderr while a command runs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initLogging(verbose, cmd.ErrOrStderr())
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Log operation details to stderr (skipped paths, check progress, etc.)")
	rootCmd.SetHelpFunc(renderHelp)
}

func hasVisibleSubcommands(cmd *cobra.Command) bool {
	for _, sub := range cmd.Commands() {
		if sub.IsAvailableCommand() && !sub.Hidden {
			return true
		}
	}
	return false
}

func renderHelp(cmd *cobra.Command, args []string) {
	if !hasVisibleSubcommands(cmd) {
		ui.RenderCommandHelp(cmd.OutOrStdout(), cmd)
		return
	}

	commands := make([]ui.HelpCommand, 0, len(cmd.Commands()))
	for _, sub := range cmd.Commands() {
		if !sub.IsAvailableCommand() || sub.Hidden {
			continue
		}
		commands = append(commands, ui.HelpCommand{
			Name:        sub.Name(),
			Description: sub.Short,
		})
	}

	ui.RenderHelp(cmd.OutOrStdout(), cmd.Short, ui.CommandDescription(cmd), commands, cmd.PersistentFlags())
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
