package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type likeOpts struct {
	sc *spotify.Client
}

func NewCmdLike(f *factory.Factory) *cobra.Command {
	opts := &likeOpts{}

	return &cobra.Command{
		Use:     "like",
		Aliases: []string{"lk"},
		Short:   "Save or unsave the currently playing track",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdLike(cmd.Context(), opts)
		},
	}
}

func runCmdLike(ctx context.Context, opts *likeOpts) error {
	pb, err := opts.sc.GetCurrentPlayback(ctx)
	if err != nil {
		return err
	}
	if pb.Item == nil {
		ui.Info("No track currently playing").Show()
		return nil
	}
	track := pb.Item
	saved, err := opts.sc.IsTrackSaved(ctx, track.ID)
	if err != nil {
		return err
	}
	if saved {
		if err := opts.sc.RemoveSavedTrack(ctx, track.ID); err != nil {
			return err
		}
		ui.Success("Unliked %s", track.Name).Show()
	} else {
		if err := opts.sc.SaveTrack(ctx, track.ID); err != nil {
			return err
		}
		ui.Success("Liked %s", track.Name).Show()
	}
	return nil
}
