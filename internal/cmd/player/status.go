package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type statusOpts struct {
	sc *spotify.Client
}

func NewCmdStatus(f *factory.Factory) *cobra.Command {
	opts := &statusOpts{}

	return &cobra.Command{
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
}

func runCmdStatus(ctx context.Context, opts *statusOpts) error {
	if _, err := spotify.GetDevice(ctx, opts.sc); err != nil {
		return err
	}
	pb, err := opts.sc.GetCurrentPlayback(ctx)
	if err != nil {
		return err
	}
	if pb.Item == nil {
		ui.Info("No active playback session").Show()
		return nil
	}
	track := pb.Item
	progressMS := pb.ProgressMS
	durationMS := track.DurationMs

	pct := 0
	if durationMS > 0 {
		pct = int(float64(progressMS) / float64(durationMS) * 100)
	}

	filled := func(s string) *ui.Style { return ui.Text(s).Yellow() }
	empty := func(s string) *ui.Style { return ui.Text(s).Dimmed() }
	state := ui.Text("paused").Yellow()
	if pb.IsPlaying {
		filled = func(s string) *ui.Style { return ui.Text(s).Green() }
		state = ui.Text("playing").Green()
	}
	shuffle := ui.Text("off").Dimmed()
	if pb.ShuffleOn {
		shuffle = ui.Text("on").Green()
	}
	repeat := ui.Text(pb.RepeatOn).Dimmed()
	if pb.RepeatOn != "off" && pb.RepeatOn != "" {
		repeat = ui.Text(pb.RepeatOn).Green()
	}
	ui.StatusRow("now playing", ui.Text(track.Name).Bold())
	ui.StatusRow("artist", ui.Text(spotify.JoinArtists(track.Artists)))
	ui.StatusRow("album", ui.Text(track.Album.Name))
	ui.StatusRow("device", ui.Text(pb.Device.Name), ui.Text("  •  vol %d%%", pb.Device.VolumePercent))

	ui.StatusRow("progress",
		ui.Text("%s / %s  ", spotify.FormatMS(progressMS), spotify.FormatMS(durationMS)),
		ui.ProgressBar(progressMS, durationMS, filled, empty),
		ui.Text("  %d%%", pct).Dimmed(),
	)
	ui.StatusRow("status", state, ui.Text("  •  shuffle "), shuffle, ui.Text("  •  repeat "), repeat)
	return nil
}
