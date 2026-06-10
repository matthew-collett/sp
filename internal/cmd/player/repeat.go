package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type repeatOpts struct {
	sc *spotify.Client
}

func NewCmdRepeat(f *factory.Factory) *cobra.Command {
	opts := &repeatOpts{}

	return &cobra.Command{
		Use:     "repeat",
		Aliases: []string{"rp"},
		Short:   "Cycle repeat mode (off → context → track)",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdRepeat(cmd.Context(), opts)
		},
	}
}

func runCmdRepeat(ctx context.Context, opts *repeatOpts) error {
	device, err := cmdutil.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}
	pb, err := opts.sc.GetCurrentPlayback(ctx)
	if err != nil {
		return err
	}
	next := nextRepeatState(pb.RepeatOn)
	if err := opts.sc.SetRepeat(ctx, device.ID, next); err != nil {
		return err
	}
	ui.Success("Repeat %s", next).Show()
	return nil
}

func nextRepeatState(current string) string {
	switch current {
	case "off":
		return "context"
	case "context":
		return "track"
	default:
		return "off"
	}
}
