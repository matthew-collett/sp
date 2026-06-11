package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type pauseOpts struct {
	sc *spotify.Client
}

func NewCmdPause(f *factory.Factory) *cobra.Command {
	opts := &pauseOpts{}
	cmd := &cobra.Command{
		Use:     "pause",
		Aliases: []string{"pa"},
		Short:   "Pause Spotify playback",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdPause(cmd.Context(), opts)
		},
	}
	return cmd
}

func runCmdPause(ctx context.Context, opts *pauseOpts) error {
	device, err := spotify.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}
	if err := opts.sc.Pause(ctx, device.ID); err != nil {
		return err
	}
	ui.Success("Paused").Show()
	return nil
}
