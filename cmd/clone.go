package cmd

import (
	"time"

	"github.com/kloudyuk/gitter/pkg/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloneCmd())
}

func cloneCmd() *cobra.Command {
	flags := struct {
		interval time.Duration
		timeout  time.Duration
		width    int
	}{}
	cmd := &cobra.Command{
		Use:   "clone URL",
		Short: "Clone a git repo repeatedly to check stability",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ui.Start(args[0], flags.interval, flags.timeout, flags.width)
		},
	}
	cmd.Flags().DurationVarP(&flags.interval, "interval", "i", 2*time.Second, "interval between clones")
	cmd.Flags().DurationVarP(&flags.timeout, "timeout", "t", 10*time.Second, "git clone timeout")
	cmd.Flags().IntVarP(&flags.width, "width", "w", 120, "terminal width for display")
	return cmd
}
