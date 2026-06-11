package queue

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

type queueAddOpts struct {
	sc    *spotify.Client
	shelf *config.Shelf
	arg   string
}

func NewCmdQueueAdd(f *factory.Factory) *cobra.Command {
	opts := &queueAddOpts{}

	cmd := &cobra.Command{
		Use:               "add <shelf-name|uri>",
		Short:             "Add a track to the playback queue",
		Args:              cobra.ExactArgs(1),
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
			opts.arg = args[0]
			return runCmdQueueAdd(cmd.Context(), opts)
		},
	}

	return cmd
}

func runCmdQueueAdd(ctx context.Context, opts *queueAddOpts) error {
	if !spotify.ValidURI(opts.arg) && !opts.shelf.Has(opts.arg) {
		return fmt.Errorf("%q is not a shelf item or a valid Spotify URI", opts.arg)
	}

	uri := opts.arg
	var name string
	if item, ok := opts.shelf.Get(opts.arg); ok {
		uri, name = item.URI, item.Name
	}

	if !spotify.IsTrack(uri) {
		return fmt.Errorf("only tracks can be added to the queue")
	}

	if err := opts.sc.AddToQueue(ctx, uri); err != nil {
		return err
	}

	if name != "" {
		ui.Success("Added %q to queue", name).Show()
	} else {
		ui.Success("Added to queue").Show()
	}
	return nil
}
