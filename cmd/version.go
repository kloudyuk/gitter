package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.0.0" // set at compile time with -ldflags "-X kloudy/gitter/cmd.Version="
var SHA = "dev"       // set at compile time with -ldflags "-X kloudy/gitter/cmd.SHA=$(git rev-parse --short HEAD)"

func init() {
	rootCmd.AddCommand(versionCmd())
}

func versionCmd() *cobra.Command {
	flags := struct {
		short bool
	}{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version info",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.short {
				fmt.Println(Version)
			} else {
				fmt.Printf("%s (%s)\n", Version, SHA)
			}
		},
	}
	cmd.Flags().BoolVarP(&flags.short, "short", "s", false, "Show short version without SHA")
	return cmd
}
