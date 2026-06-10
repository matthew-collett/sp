package cmd

import (
	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdVersion(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version for sp",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmdVersion(f)
		},
	}

	cmdutil.DisableConfigCheck(cmd)
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdVersion(f *factory.Factory) error {
	f.Version.Style().Show()
	return nil
}
