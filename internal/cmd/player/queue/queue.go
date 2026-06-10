package queue

import (
	"context"
	"fmt"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type queueOpts struct {
	sc      *spotify.Client
	noPager bool
}

func NewCmdQueue(f *factory.Factory) *cobra.Command {
	opts := &queueOpts{}

	cmd := &cobra.Command{
		Use:     "queue",
		Aliases: []string{"q"},
		Short:   "Show the current playback queue",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdQueue(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.noPager, "no-pager", false, "disable pager for output")

	cmd.AddCommand(NewCmdQueueAdd(f))

	return cmd
}

func runCmdQueue(ctx context.Context, opts *queueOpts) error {
	res, err := opts.sc.GetQueue(ctx)
	if err != nil {
		return err
	}

	t := ui.NewTable()
	t.Title(ui.Text("tracks"))
	t.Header("#", "TRACK", "ARTISTS", "ALBUM")

	if res.CurrentlyPlaying != nil {
		t.Row(
			ui.Text("now").Dimmed().Fixed(),
			ui.Text(res.CurrentlyPlaying.Name).Bold(),
			ui.Text(cmdutil.JoinArtists(res.CurrentlyPlaying.Artists)).Fixed(),
			ui.Text(res.CurrentlyPlaying.Album.Name).Dimmed().Fixed(),
		)
	}

	for i, track := range res.Queue {
		t.Row(
			ui.Text(fmt.Sprintf("%d", i+1)).Dimmed().Fixed(),
			ui.Text(track.Name).Bold(),
			ui.Text(cmdutil.JoinArtists(track.Artists)).Fixed(),
			ui.Text(track.Album.Name).Dimmed().Fixed(),
		)
	}

	if t.Empty() {
		ui.Info("Queue is empty").Show()
		return nil
	}

	t.Render(!opts.noPager)
	return nil
}
