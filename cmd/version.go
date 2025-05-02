package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func init() {
	rootCmd.AddCommand(versionCmd())
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
	return cmd
}
