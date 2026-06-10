package devices

import (
	"errors"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/osutil"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

func NewCmdOpen(_ *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open Spotify on this device",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmdOpen()
		},
	}

	cmdutil.DisableConfigCheck(cmd)
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdOpen() error {
	if err := osutil.OpenApplication("Spotify"); err != nil {
		return errors.New("failed to open spotify")
	}
	ui.Success("Spotify opened").Show()
	return nil
}
