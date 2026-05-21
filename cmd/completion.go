package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for eggl.

To load completions:

Bash:
  $ source <(eggl completion bash)
  # or persist:
  $ eggl completion bash > /etc/bash_completion.d/eggl

Zsh:
  $ source <(eggl completion zsh)
  # or persist:
  $ eggl completion zsh > "${fpath[1]}/_eggl"

Fish:
  $ eggl completion fish | source
  # or persist:
  $ eggl completion fish > ~/.config/fish/completions/eggl.fish

PowerShell:
  PS> eggl completion powershell | Out-String | Invoke-Expression
  # or persist:
  PS> eggl completion powershell > eggl.ps1
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
