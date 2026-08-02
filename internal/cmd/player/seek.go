package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type seekOpts struct {
	sc       *spotify.Client
	position string
}

func NewCmdSeek(f *factory.Factory) *cobra.Command {
	opts := &seekOpts{}

	return &cobra.Command{
		Use:   "seek <position>",
		Short: "Seek to a position in the current track (e.g. 1:30 or 90)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			opts.position = args[0]
			return runCmdSeek(cmd.Context(), opts)
		},
	}
}

func runCmdSeek(ctx context.Context, opts *seekOpts) error {
	device, err := spotify.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}
	ms, err := spotify.ParsePosition(opts.position)
	if err != nil {
		return err
	}

	name := ""
	if pb, err := opts.sc.GetCurrentPlayback(ctx); err == nil && pb.Item != nil {
		name = pb.Item.Name
	}

	if err := opts.sc.Seek(ctx, device.ID, ms); err != nil {
		return err
	}

	if name != "" {
		ui.Success("Seeked to %s in %q", opts.position, name).Show()
	} else {
		ui.Success("Seeked to %s", opts.position).Show()
	}
	return nil
}
