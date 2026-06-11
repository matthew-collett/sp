package player

import (
	"context"
	"fmt"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type playOpts struct {
	sc       *spotify.Client
	shelf    *config.Shelf
	arg      string
	position int
	shuffle  bool
}

func NewCmdPlay(f *factory.Factory) *cobra.Command {
	opts := &playOpts{}

	cmd := &cobra.Command{
		Use:               "play [shelf-name|uri]",
		Aliases:           []string{"pl"},
		Short:             "Start or resume Spotify playback",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cmdutil.ShelfCompletion(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			shelf, err := f.Shelf()
			if err != nil {
				return err
			}
			opts.sc = sc
			opts.shelf = shelf
			if len(args) > 0 {
				opts.arg = args[0]
			}
			return runCmdPlay(cmd.Context(), opts)
		},
	}

	cmd.Flags().IntVarP(&opts.position, "position", "p", 0, "start position in seconds")
	cmd.Flags().BoolVarP(&opts.shuffle, "shuffle", "s", false, "enable shuffle for albums and playlists")

	return cmd
}

func runCmdPlay(ctx context.Context, opts *playOpts) error {
	device, err := spotify.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}

	var req spotify.PlayPlaybackRequest
	if opts.position > 0 {
		req.PositionMS = opts.position * 1000
	}

	var name string
	if opts.arg != "" {
		if !spotify.ValidURI(opts.arg) && !opts.shelf.Has(opts.arg) {
			return fmt.Errorf("%q is not a shelf item or a valid Spotify URI", opts.arg)
		}
		uri := opts.arg
		if item, ok := opts.shelf.Get(opts.arg); ok {
			uri, name = item.URI, item.Name
		}
		if spotify.IsTrack(uri) {
			req.URIs = []string{uri}
		} else {
			req.ContextURI = uri
		}
	}

	if err := opts.sc.Play(ctx, device.ID, req); err != nil {
		return err
	}

	if opts.shuffle {
		if err := opts.sc.SetShuffle(ctx, device.ID, true); err != nil {
			return err
		}
	}

	if name != "" {
		ui.Success("Playing %s", name).Show()
	} else {
		ui.Success("Playing").Show()
	}
	return nil
}
