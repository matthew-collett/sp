package devices

import (
	"context"
	"fmt"
	"strings"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type activateOpts struct {
	sc         *spotify.Client
	deviceName string
}

func NewCmdActivate(f *factory.Factory) *cobra.Command {
	opts := &activateOpts{}

	cmd := &cobra.Command{
		Use:     "activate <name>",
		Aliases: []string{"act"},
		Short:   "Activate a device for playback by its name",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			opts.deviceName = args[0]
			return runCmdActivate(cmd.Context(), opts)
		},
	}

	return cmd
}

func runCmdActivate(ctx context.Context, opts *activateOpts) error {
	resp, err := opts.sc.GetDevices(ctx)
	if err != nil {
		return err
	}

	if len(resp.Devices) == 0 {
		return spotify.ErrNoAvailableDevices
	}

	var device *spotify.Device
	for _, d := range resp.Devices {
		if strings.EqualFold(d.Name, opts.deviceName) {
			device = &d
			break
		}
	}

	if device == nil {
		return fmt.Errorf("device %q not found", opts.deviceName)
	}

	if !device.IsActive {
		req := spotify.TransferPlaybackRequest{
			DeviceIDs: []string{device.ID},
		}
		if err := opts.sc.TransferPlayback(ctx, req); err != nil {
			return err
		}
	}

	ui.Success("Activated %q", device.Name).Show()
	return nil
}
