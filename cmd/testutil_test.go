package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func runCmd(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	resetCommandFlags(rootCmd)
	rootCmd.SetArgs(args)

	var out, errOut bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&errOut)

	err = rootCmd.Execute()
	return out.String(), errOut.String(), err
}

func runHelp(t *testing.T, args ...string) string {
	t.Helper()

	stdout, _, err := runCmd(t, args...)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	return stdout
}

func resetCommandFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetCommandFlags(sub)
	}
}
