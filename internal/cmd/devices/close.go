package devices

import (
	"errors"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/osutil"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

func NewCmdClose(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close Spotify on this device",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmdClose()
		},
	}

	return cmd
}

func runCmdClose() error {
	if err := osutil.CloseApplication("Spotify"); err != nil {
		return errors.New("failed to close spotify")
	}

	ui.Success("Spotify closed").Show()
	return nil
}
