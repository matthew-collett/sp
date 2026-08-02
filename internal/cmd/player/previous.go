package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type previousOpts struct {
	sc *spotify.Client
}

func NewCmdPrevious(f *factory.Factory) *cobra.Command {
	opts := &previousOpts{}

	cmd := &cobra.Command{
		Use:     "previous",
		Aliases: []string{"prev"},
		Short:   "Skip to previous track",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdPrevious(cmd.Context(), opts)
		},
	}

	return cmd
}

func runCmdPrevious(ctx context.Context, opts *previousOpts) error {
	device, err := spotify.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}

	if err := opts.sc.Previous(ctx, device.ID); err != nil {
		return err
	}

	name := ""
	if pb, err := spotify.GetCurrentPlaybackDelayed(ctx, opts.sc); err == nil && pb.Item != nil {
		name = pb.Item.Name
	}
	if name != "" {
		ui.Success("Previous track: %q", name).Show()
	} else {
		ui.Success("Previous track").Show()
	}
	return nil
}
