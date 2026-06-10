package player

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/matthew-collett/sp/internal/cmdutil"
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
	device, err := cmdutil.GetDevice(ctx, opts.sc)
	if err != nil {
		return err
	}
	ms, err := parsePosition(opts.position)
	if err != nil {
		return err
	}
	if err := opts.sc.Seek(ctx, device.ID, ms); err != nil {
		return err
	}
	ui.Success("Seeked to %s", opts.position).Show()
	return nil
}

func parsePosition(s string) (int, error) {
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		mins, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid position %q", s)
		}
		secs, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid position %q", s)
		}
		return (mins*60 + secs) * 1000, nil
	}
	secs, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid position %q — use m:ss or seconds", s)
	}
	return secs * 1000, nil
}
