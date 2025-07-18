package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	// Version can be set via ldflags: -ldflags "-X github.com/kloudyuk/gitter/cmd.Version=v1.2.3"
	Version = "dev"
)

// getVersion returns the version string, trying multiple sources
func getVersion() string {
	// If Version was set via ldflags and it's not the default "dev", use it
	if Version != "dev" {
		return Version
	}

	// Try to get version from build info (works with go install)
	if info, ok := debug.ReadBuildInfo(); ok {
		// Check the module version (this works with go install github.com/user/repo@version)
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}

		// Check VCS information for tag or revision
		var revision string
		for _, setting := range info.Settings {
			if setting.Key == "vcs.tag" && setting.Value != "" {
				return setting.Value
			}
			if setting.Key == "vcs.revision" && len(setting.Value) >= 7 {
				revision = setting.Value[:7]
			}
		}

		// If we have a revision but no tag, return it
		if revision != "" {
			return revision
		}
	}

	// Fall back to the default
	return Version
}

func init() {
	rootCmd.AddCommand(versionCmd())
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}
	return cmd
}
