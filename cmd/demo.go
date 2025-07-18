package cmd

import (
	"time"

	"github.com/kloudyuk/gitter/pkg/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(demoCmd())
}

func demoCmd() *cobra.Command {
	flags := struct {
		interval time.Duration
		timeout  time.Duration
		width    int
	}{}
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Run in demo mode with simulated git operations",
		Long: `Demo mode simulates git clone operations without actually cloning repositories.
This is useful for testing the UI, demonstrating features, and development.
It will show simulated successes and failures with realistic timing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ui.StartDemo(flags.interval, flags.timeout, flags.width)
		},
	}
	cmd.Flags().DurationVarP(&flags.interval, "interval", "i", 2*time.Second, "interval between simulated clones")
	cmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 10*time.Second, "timeout for demo operations")
	cmd.Flags().IntVarP(&flags.width, "width", "w", 100, "terminal width for display")
	return cmd
}
