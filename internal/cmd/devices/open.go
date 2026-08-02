package devices

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/osutil"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

const (
	openPollInterval = 500 * time.Millisecond
	openPollTimeout  = 15 * time.Second
)

type openOpts struct {
	sc         *spotify.Client
	deviceName string
}

func NewCmdOpen(f *factory.Factory) *cobra.Command {
	opts := &openOpts{}

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open Spotify on this device, optionally activating a device",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.deviceName != "" {
				sc, err := f.SpotifyClient(cmd.Context())
				if err != nil {
					return err
				}
				opts.sc = sc
			}
			return runCmdOpen(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.deviceName, "activate", "a", "", "activate this device after opening")

	cmdutil.DisableConfigCheck(cmd)
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdOpen(ctx context.Context, opts *openOpts) error {
	if err := osutil.OpenApplication("Spotify"); err != nil {
		return errors.New("failed to open spotify")
	}

	if opts.deviceName == "" {
		ui.Success("Spotify opened").Show()
		return nil
	}

	device, err := waitForDevice(ctx, opts.sc, opts.deviceName)
	if err != nil {
		return err
	}

	if !device.IsActive {
		req := spotify.TransferPlaybackRequest{DeviceIDs: []string{device.ID}}
		if err := opts.sc.TransferPlayback(ctx, req); err != nil {
			return err
		}
	}

	ui.Success("Opened and activated %q", device.Name).Show()
	return nil
}

func waitForDevice(ctx context.Context, sc *spotify.Client, name string) (*spotify.Device, error) {
	defer ui.ClearLine()

	ctx, cancel := context.WithTimeout(ctx, openPollTimeout)
	defer cancel()

	for tick := 0; ; tick++ {
		resp, err := sc.GetDevices(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, fmt.Errorf("device %q did not appear within %s", name, openPollTimeout)
			}
			return nil, err
		}
		for _, d := range resp.Devices {
			if strings.EqualFold(d.Name, name) {
				return &d, nil
			}
		}

		ui.Spinner(tick, "Waiting for %q to appear", name)

		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, fmt.Errorf("device %q did not appear within %s", name, openPollTimeout)
			}
			return nil, ctx.Err()
		case <-time.After(openPollInterval):
		}
	}
}
