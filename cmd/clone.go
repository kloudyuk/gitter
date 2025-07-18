package cmd

import (
	"fmt"
	"time"

	"github.com/kloudyuk/gitter/pkg/ui"

	"github.com/spf13/cobra"
)

const (
	// Input validation constants
	MinWidth = 50
	MaxWidth = 300
)

func init() {
	rootCmd.AddCommand(cloneCmd())
}

func cloneCmd() *cobra.Command {
	flags := struct {
		interval    time.Duration
		timeout     time.Duration
		width       int
		demo        bool
		errorHistory int
	}{}
	cmd := &cobra.Command{
		Use:   "clone URL",
		Short: "Clone a git repo repeatedly to check stability",
		Long: `Clone a git repository repeatedly to test its stability and reliability.
Use the --demo flag to run in simulation mode without actually cloning repositories.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate input parameters
			if flags.interval <= 0 {
				return fmt.Errorf("interval must be positive, got %v", flags.interval)
			}
			if flags.timeout <= 0 {
				return fmt.Errorf("timeout must be positive, got %v", flags.timeout)
			}
			if flags.width < MinWidth || flags.width > MaxWidth {
				return fmt.Errorf("width must be between %d and %d, got %d", MinWidth, MaxWidth, flags.width)
			}
			if flags.errorHistory <= 0 {
				return fmt.Errorf("error-history must be positive, got %d", flags.errorHistory)
			}

			var repoURL string
			if flags.demo {
				repoURL = "https://github.com/demo/repo.git (simulated)"
			} else {
				if len(args) == 0 {
					return fmt.Errorf("repository URL is required when not in demo mode")
				}
				repoURL = args[0]
			}
			return ui.Start(repoURL, flags.interval, flags.timeout, flags.width, flags.demo, flags.errorHistory)
		},
	}
	cmd.Flags().DurationVarP(&flags.interval, "interval", "i", 2*time.Second, "interval between clones (must be positive)")
	cmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 10*time.Second, "timeout for clone operations (must be positive)")
	cmd.Flags().IntVarP(&flags.width, "width", "w", 100, fmt.Sprintf("terminal width for display (%d-%d)", MinWidth, MaxWidth))
	cmd.Flags().BoolVarP(&flags.demo, "demo", "d", false, "run in demo mode with simulated git operations")
	cmd.Flags().IntVarP(&flags.errorHistory, "error-history", "e", 5, "number of recent errors to display (must be positive)")
	return cmd
}
