package devices

import (
	"context"
	"sort"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type devicesOpts struct {
	sc      *spotify.Client
	user    string
	noPager bool
}

func NewCmdDevices(f *factory.Factory) *cobra.Command {
	opts := &devicesOpts{}

	cmd := &cobra.Command{
		Use:     "devices",
		Aliases: []string{"ds"},
		Short:   "List available Spotify devices",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := f.SpotifyClient(cmd.Context())
			if err != nil {
				return err
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.sc = sc
			opts.user = cfg.User
			return runCmdDevices(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.noPager, "no-pager", false, "disable pager for output")

	return cmd
}

func runCmdDevices(ctx context.Context, opts *devicesOpts) error {
	resp, err := opts.sc.GetDevices(ctx)
	if err != nil {
		return err
	}

	if len(resp.Devices) == 0 {
		return cmdutil.ErrNoAvailableDevices
	}

	sort.Slice(resp.Devices, func(i, j int) bool {
		return resp.Devices[i].IsActive && !resp.Devices[j].IsActive
	})

	t := ui.NewTable()
	t.Title(ui.Text("devices for %s", ui.Text(opts.user).Bold()))
	t.Header("NAME", "TYPE", "VOLUME", "ACTIVE", "PRIVATE", "RESTRICTED")

	for _, device := range resp.Devices {
		t.Row(
			ui.Text(device.Name).Bold(),
			ui.Text(device.Type).Bold(),
			ui.Text("%d%%", device.VolumePercent).Dimmed(),
			ui.Bool(device.IsActive),
			ui.Bool(device.IsPrivateSession),
			ui.Bool(device.IsRestricted),
		)
	}

	t.Render(!opts.noPager)
	return nil
}
