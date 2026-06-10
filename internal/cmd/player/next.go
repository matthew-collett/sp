package player

import (
	"context"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type nextOpts struct {
	sc *spotify.Client
}

func NewCmdNext(f *factory.Factory) *cobra.Command {
	opts := &nextOpts{}

	cmd := &cobra.Command{
		Use:     "next",
		Aliases: []string{"n"},
		Short:   "Skip to next track",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			opts.sc = sc
			return runCmdNext(cmd.Context(), opts)
		},
	}

	return cmd
}

func runCmdNext(ctx context.Context, opts *nextOpts) error {
	device, err := cmdutil.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}

	if err := opts.sc.Next(ctx, device.ID); err != nil {
		return err
	}
	ui.Success("Next track").Show()
	return nil
}
