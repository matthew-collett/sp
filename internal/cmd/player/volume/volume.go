package volume

import (
	"context"
	"fmt"
	"strconv"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type volumeOpts struct {
	sc     *spotify.Client
	volume string
}

func NewCmdVolume(f *factory.Factory) *cobra.Command {
	opts := &volumeOpts{}

	cmd := &cobra.Command{
		Use:     "volume [level]",
		Aliases: []string{"v", "vol"},
		Short:   "Show current volume or set volume level (0-100)",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			if len(args) > 0 {
				opts.volume = args[0]
			}
			return runCmdVolume(cmd.Context(), opts)
		},
	}

	cmd.AddCommand(
		NewCmdVolumeUp(f),
		NewCmdVolumeDown(f),
	)

	return cmd
}

func runCmdVolume(ctx context.Context, opts *volumeOpts) error {
	state, err := opts.sc.GetPlaybackState(ctx)
	if err != nil {
		return err
	}

	if state.Device.ID == "" {
		return cmdutil.ErrNoActiveDevice
	}

	if opts.volume == "" {
		ui.Text("Volume: %d%%", state.Device.VolumePercent).Show()
		return nil
	}

	v, err := strconv.Atoi(opts.volume)
	if err != nil {
		return fmt.Errorf("invalid volume level: %s (must be 0-100)", opts.volume)
	}
	if v < 0 || v > 100 {
		return fmt.Errorf("invalid volume level: %d (must be between 0 and 100)", v)
	}

	if err := opts.sc.SetVolume(ctx, state.Device.ID, v); err != nil {
		return err
	}

	ui.Success("Volume set to %d%%", v).Show()
	return nil
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
