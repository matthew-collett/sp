package volume

import (
	"context"
	"fmt"
	"strconv"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type volumeDownOpts struct {
	sc    *spotify.Client
	delta int
}

func NewCmdVolumeDown(f *factory.Factory) *cobra.Command {
	opts := &volumeDownOpts{}

	cmd := &cobra.Command{
		Use:   "down [N]",
		Short: "Decrease volume by N (default 10)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			opts.delta = 10
			if len(args) > 0 {
				n, err := strconv.Atoi(args[0])
				if err != nil || n < 0 {
					return fmt.Errorf("invalid amount: %s (must be a positive number)", args[0])
				}
				opts.delta = n
			}
			return runCmdVolumeDown(cmd.Context(), opts)
		},
	}

	return cmd
}

func runCmdVolumeDown(ctx context.Context, opts *volumeDownOpts) error {
	state, err := opts.sc.GetPlaybackState(ctx)
	if err != nil {
		return err
	}

	if state.Device.ID == "" {
		return spotify.ErrNoActiveDevice
	}

	v := clamp(state.Device.VolumePercent-opts.delta, 0, 100)
	if err := opts.sc.SetVolume(ctx, state.Device.ID, v); err != nil {
		return err
	}

	ui.Success("Volume set to %d%%", v).Show()
	return nil
}
