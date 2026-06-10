package search

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

func parseType(s string) (string, bool) {
	switch strings.ToLower(s) {
	case "tracks":
		return "track", true
	case "albums":
		return "album", true
	case "playlists":
		return "playlist", true
	case "artists":
		return "artist", true
	}
	return "", false
}

type searchOpts struct {
	sc      *spotify.Client
	query   string
	typ     string
	artist  string
	mine    bool
	limit   int
	noPager bool
}

func NewCmdSearch(f *factory.Factory) *cobra.Command {
	opts := &searchOpts{}

	cmd := &cobra.Command{
		Use:     "search <tracks|albums|playlists|artists> [query]",
		Aliases: []string{"s"},
		Short:   "Search for tracks, albums, playlists, or artists",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			typ, ok := parseType(args[0])
			if !ok {
				return fmt.Errorf("invalid type %q — must be one of: tracks, albums, playlists, artists", args[0])
			}
			opts.typ = typ
			if len(args) == 2 {
				opts.query = args[1]
			} else if !opts.mine {
				return errors.New("query required unless --mine is set")
			}
			if opts.mine && cmd.Flags().Changed("limit") {
				return errors.New("--limit cannot be used with --mine")
			}
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdSearch(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.artist, "artist", "a", "", "filter by artist name")
	cmd.Flags().BoolVarP(&opts.mine, "mine", "m", false, "search only your library")
	cmd.Flags().IntVarP(&opts.limit, "limit", "l", 50, "number of results to show")
	cmd.Flags().BoolVar(&opts.noPager, "no-pager", false, "disable pager for output")

	return cmd
}

func runCmdSearch(ctx context.Context, opts *searchOpts) error {
	if opts.mine {
		return searchLibrary(ctx, opts)
	}
	return searchSpotify(ctx, opts)
}

func searchSpotify(ctx context.Context, opts *searchOpts) error {
	query := opts.query
	if opts.artist != "" {
		query = fmt.Sprintf("%s artist:%s", query, opts.artist)
	}

	res, err := opts.sc.Search(ctx, query, []string{opts.typ}, opts.limit)
	if err != nil {
		return err
	}

	t := ui.NewTable()
	t.Title(ui.Text(`%ss matching "%s"`, opts.typ, opts.query))

	switch opts.typ {
	case "track":
		t.Header("NAME", "ARTISTS", "DURATION", "URI")
		if res.Tracks != nil {
			for _, track := range res.Tracks.Items {
				t.Row(
					ui.Text(track.Name).Bold(),
					ui.Text(cmdutil.JoinArtists(track.Artists)).Dimmed(),
					ui.Text(formatDuration(track.DurationMs)).Dimmed().Fixed(),
					ui.Text(track.URI).Yellow().Fixed(),
				)
			}
		}
	case "album":
		t.Header("NAME", "ARTISTS", "TRACKS", "URI")
		if res.Albums != nil {
			for _, album := range res.Albums.Items {
				t.Row(
					ui.Text(album.Name).Bold(),
					ui.Text(cmdutil.JoinArtists(album.Artists)).Dimmed(),
					ui.Text("%d", album.TotalTracks).Dimmed().Fixed(),
					ui.Text(album.URI).Yellow().Fixed(),
				)
			}
		}
	case "playlist":
		t.Header("NAME", "OWNER", "PUBLIC", "TRACKS", "URI")
		if res.Playlists != nil {
			for _, p := range res.Playlists.Items {
				t.Row(
					ui.Text(p.Name).Bold(),
					ui.Text(p.Owner.DisplayName).Dimmed().Fixed(),
					ui.Bool(p.Public).Fixed(),
					ui.Text("%d", p.Tracks.Total).Dimmed().Fixed(),
					ui.Text(p.URI).Yellow().Fixed(),
				)
			}
		}
	case "artist":
		t.Header("NAME", "GENRES", "FOLLOWERS", "URI")
		if res.Artists != nil {
			for _, a := range res.Artists.Items {
				t.Row(
					ui.Text(a.Name).Bold(),
					ui.Text(strings.Join(a.Genres, ", ")).Dimmed(),
					ui.Text(formatNumber(a.Followers.Total)).Dimmed().Fixed(),
					ui.Text(a.URI).Yellow().Fixed(),
				)
			}
		}
	}

	if t.Empty() {
		return fmt.Errorf("no %ss found for %q", opts.typ, opts.query)
	}

	t.Render(!opts.noPager)
	return nil
}

func searchLibrary(ctx context.Context, opts *searchOpts) error {
	q := strings.ToLower(opts.query)

	t := ui.NewTable()
	if opts.query != "" {
		t.Title(ui.Text(`%ss matching "%s"`, opts.typ, opts.query))
	} else {
		t.Title(ui.Text("%ss", opts.typ))
	}

	switch opts.typ {
	case "track":
		items, err := opts.sc.GetSavedTracks(ctx)
		if err != nil {
			return err
		}
		t.Header("NAME", "ARTISTS", "DURATION", "URI")
		for _, saved := range items {
			if q != "" && !strings.Contains(strings.ToLower(saved.Track.Name), q) {
				continue
			}
			t.Row(
				ui.Text(saved.Track.Name).Bold(),
				ui.Text(cmdutil.JoinArtists(saved.Track.Artists)).Dimmed(),
				ui.Text(formatDuration(saved.Track.DurationMs)).Dimmed().Fixed(),
				ui.Text(saved.Track.URI).Yellow().Fixed(),
			)
		}
	case "album":
		items, err := opts.sc.GetSavedAlbums(ctx)
		if err != nil {
			return err
		}
		t.Header("NAME", "ARTISTS", "TRACKS", "URI")
		for _, saved := range items {
			if q != "" && !strings.Contains(strings.ToLower(saved.Album.Name), q) {
				continue
			}
			t.Row(
				ui.Text(saved.Album.Name).Bold(),
				ui.Text(cmdutil.JoinArtists(saved.Album.Artists)).Dimmed(),
				ui.Text("%d", saved.Album.TotalTracks).Dimmed().Fixed(),
				ui.Text(saved.Album.URI).Yellow().Fixed(),
			)
		}
	case "playlist":
		items, err := opts.sc.GetPlaylists(ctx)
		if err != nil {
			return err
		}
		t.Header("NAME", "OWNER", "PUBLIC", "TRACKS", "URI")
		for _, p := range items {
			if q != "" && !strings.Contains(strings.ToLower(p.Name), q) {
				continue
			}
			t.Row(
				ui.Text(p.Name).Bold(),
				ui.Text(p.Owner.DisplayName).Dimmed().Fixed(),
				ui.Bool(p.Public).Fixed(),
				ui.Text("%d", p.Tracks.Total).Dimmed().Fixed(),
				ui.Text(p.URI).Yellow().Fixed(),
			)
		}
	case "artist":
		items, err := opts.sc.GetFollowedArtists(ctx)
		if err != nil {
			return err
		}
		t.Header("NAME", "GENRES", "FOLLOWERS", "URI")
		for _, a := range items {
			if q != "" && !strings.Contains(strings.ToLower(a.Name), q) {
				continue
			}
			t.Row(
				ui.Text(a.Name).Bold(),
				ui.Text(strings.Join(a.Genres, ", ")).Dimmed(),
				ui.Text(formatNumber(a.Followers.Total)).Dimmed().Fixed(),
				ui.Text(a.URI).Yellow().Fixed(),
			)
		}
	}

	if t.Empty() {
		if opts.query != "" {
			return fmt.Errorf("no %ss found in your library for %q", opts.typ, opts.query)
		}
		return fmt.Errorf("no %ss found in your library", opts.typ)
	}

	t.Render(!opts.noPager)
	return nil
}

func formatDuration(ms int) string {
	total := ms / 1000
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}

func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	result := make([]byte, 0, len(s)+(len(s)-1)/3)
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
