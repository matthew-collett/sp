package player

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

const (
	statusWatchInterval = time.Second
	artFetchTimeout     = 3 * time.Second
	artCols             = 16
	artRows             = 8
)

type statusOpts struct {
	sc    *spotify.Client
	watch bool
}

func NewCmdStatus(f *factory.Factory) *cobra.Command {
	opts := &statusOpts{}

	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"st"},
		Short:   "Show the current playback status",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdStatus(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.watch, "watch", "w", false, "refresh continuously")

	return cmd
}

func runCmdStatus(ctx context.Context, opts *statusOpts) error {
	if !opts.watch {
		panel, err := statusPanel(ctx, opts.sc, make(map[string]string))
		if err != nil {
			return err
		}
		ui.Print(panel)
		return nil
	}
	return watch(ctx, opts.sc)
}

func statusPanel(ctx context.Context, sc *spotify.Client, genres map[string]string) (string, error) {
	if _, err := spotify.GetDevice(ctx, sc); err != nil {
		return "", err
	}
	pb, err := sc.GetCurrentPlayback(ctx)
	if err != nil {
		return "", err
	}
	if pb.Item == nil {
		return ui.Info("No active playback session").String() + "\n", nil
	}

	box := ui.NewBox("Playback")
	textWidth := box.Width() - artCols - 2
	genre := genreFor(ctx, sc, pb.Item.Artists, genres)
	lines := trackLines(pb, genre, textWidth)

	for _, row := range ui.SideBySide(albumArt(ctx, pb.Item), lines, artCols) {
		box.Row(row)
	}

	box.Row("")
	label := spotify.FormatMS(pb.ProgressMS) + "/" + spotify.FormatMS(pb.Item.DurationMs)
	box.Row(ui.ProgressBarOverlay(box.Width(), progressRatio(pb.ProgressMS, pb.Item.DurationMs), label))

	return box.Render(), nil
}

func albumArt(ctx context.Context, track *spotify.Track) []string {
	blank := make([]string, artRows/2)
	for i := range blank {
		blank[i] = strings.Repeat(" ", artCols)
	}
	if len(track.Album.Images) == 0 {
		return blank
	}
	ctx, cancel := context.WithTimeout(ctx, artFetchTimeout)
	defer cancel()
	art, err := ui.RenderImage(ctx, track.Album.Images[0].URL, artCols, artRows)
	if err != nil {
		return blank
	}
	return art
}

func progressRatio(progressMS, durationMS int) float64 {
	if durationMS == 0 {
		return 0
	}
	return float64(progressMS) / float64(durationMS)
}

func trackLines(pb *spotify.CurrentPlayback, genre string, width int) []string {
	track := pb.Item

	icon := ui.Pause()
	if pb.IsPlaying {
		icon = ui.Play()
	}

	main, artists := ui.Fit(fmt.Sprintf("%s %s", icon, track.Name), spotify.JoinArtists(track.Artists), width)
	title := ui.Text("%s", main).Bold().String()
	if artists != "" {
		title += ui.Text(" %s %s", ui.Dot(), artists).String()
	}

	if genre == "" {
		genre = "no genre"
	}
	album, genreText := ui.Fit(track.Album.Name, genre, width)
	albumLine := ui.Text("%s", album).String()
	if genreText != "" {
		albumLine += ui.Text(" %s %s", ui.Dot(), genreText).Italic().Dimmed().String()
	}

	shuffle := "off"
	if pb.ShuffleOn {
		shuffle = "on"
	}
	meta := ui.Truncate(fmt.Sprintf("repeat: %s | shuffle: %s | volume: %d%% | device: %s",
		pb.RepeatOn, shuffle, pb.Device.VolumePercent, pb.Device.Name), width)

	return []string{title, albumLine, ui.Text("%s", meta).Dimmed().String()}
}

func genreFor(ctx context.Context, sc *spotify.Client, artists []spotify.Artist, cache map[string]string) string {
	if len(artists) == 0 {
		return ""
	}
	id := artists[0].ID
	if genre, ok := cache[id]; ok {
		return genre
	}
	genre := ""
	if artist, err := sc.GetArtist(ctx, id); err == nil && len(artist.Genres) > 0 {
		genre = artist.Genres[0]
	}
	cache[id] = genre
	return genre
}

func watch(ctx context.Context, sc *spotify.Client) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	genres := make(map[string]string)
	ticker := time.NewTicker(statusWatchInterval)
	defer ticker.Stop()

	for {
		panel, err := statusPanel(ctx, sc, genres)

		ui.ClearScreen()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			ui.Error(err).Show()
		} else {
			ui.Print(panel)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}
