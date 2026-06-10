package cmd

import (
	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

func NewCmdSplash(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "splash",
		Short: "Show the sp splash",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmdSplash()
		},
	}

	cmdutil.DisableConfigCheck(cmd)
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdSplash() error {
	ui.ShowSplash()
	return nil
}
