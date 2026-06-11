package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type shuffleOpts struct {
	sc *spotify.Client
}

func NewCmdShuffle(f *factory.Factory) *cobra.Command {
	opts := &shuffleOpts{}

	return &cobra.Command{
		Use:     "shuffle",
		Aliases: []string{"sh"},
		Short:   "Toggle shuffle for the current playback",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdShuffle(cmd.Context(), opts)
		},
	}
}

func runCmdShuffle(ctx context.Context, opts *shuffleOpts) error {
	device, err := spotify.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}
	pb, err := opts.sc.GetCurrentPlayback(ctx)
	if err != nil {
		return err
	}
	next := !pb.ShuffleOn
	if err := opts.sc.SetShuffle(ctx, device.ID, next); err != nil {
		return err
	}
	if next {
		ui.Success("Shuffle on").Show()
	} else {
		ui.Success("Shuffle off").Show()
	}
	return nil
}
