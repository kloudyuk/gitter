package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var red = color.New(color.FgRed)

var rootCmd = &cobra.Command{
	Use:           "gitter",
	Short:         "Gitter is a simple utility for testing a git server",
	SilenceErrors: true,
	SilenceUsage:  true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		red.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
