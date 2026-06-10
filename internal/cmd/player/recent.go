package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type recentOpts struct {
	sc      *spotify.Client
	noPager bool
}

func NewCmdRecent(f *factory.Factory) *cobra.Command {
	opts := &recentOpts{}

	cmd := &cobra.Command{
		Use:     "recent",
		Aliases: []string{"rc"},
		Short:   "Show recently played tracks",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdRecent(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.noPager, "no-pager", false, "disable pager for output")

	return cmd
}

func runCmdRecent(ctx context.Context, opts *recentOpts) error {
	res, err := opts.sc.GetRecentlyPlayed(ctx)
	if err != nil {
		return err
	}

	t := ui.NewTable()
	t.Title(ui.Text("recently played"))
	t.Header("TRACK", "ARTISTS", "ALBUM", "URI")

	for _, item := range res.Items {
		t.Row(
			ui.Text(item.Track.Name).Bold(),
			ui.Text(cmdutil.JoinArtists(item.Track.Artists)).Dimmed(),
			ui.Text(item.Track.Album.Name).Dimmed(),
			ui.Text(item.Track.URI).Yellow().Fixed(),
		)
	}

	if t.Empty() {
		ui.Info("No recently played tracks").Show()
		return nil
	}

	t.Render(!opts.noPager)
	return nil
}
