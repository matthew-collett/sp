package cmd

import (
	"github.com/MakeNowJust/heredoc"
	cmdconfig "github.com/matthew-collett/sp/internal/cmd/config"
	"github.com/matthew-collett/sp/internal/cmd/devices"
	"github.com/matthew-collett/sp/internal/cmd/player"
	"github.com/matthew-collett/sp/internal/cmd/player/queue"
	"github.com/matthew-collett/sp/internal/cmd/player/volume"
	"github.com/matthew-collett/sp/internal/cmd/search"
	"github.com/matthew-collett/sp/internal/cmd/shelf"
	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *factory.Factory) (*cobra.Command, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, config.ErrNoConfig
	}
	auth, err := f.Auth()
	if err != nil {
		return nil, config.ErrNotAuthenticated
	}

	cmd := &cobra.Command{
		Use:  "sp",
		Long: ui.Text("sp - A fast, minimal Spotify CLI for your terminal.").Bold().String(),
		Example: heredoc.Doc(`
		$ sp play
		$ sp play my-playlist
		$ sp search tracks "skinny love"
		$ sp search albums --mine
		$ sp shelf add lofi spotify:playlist:37...
		$ sp shelf
		$ sp devices
		$ sp activate macbook
		$ sp volume 50
		$ sp volume up 10
		$ sp volume down 5
		$ sp queue
		$ sp queue add spotify:track:4i...
		$ sp pause`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmdRoot(cmd)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmdutil.IsConfigCheckEnabled(cmd) && !cfg.IsConfigured() {
				return config.ErrNoConfig
			}
			if cmdutil.IsAuthCheckEnabled(cmd) && !auth.IsAuthenticated() {
				return config.ErrNotAuthenticated
			}
			return nil
		},
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.Version = f.Version.String()
	cmd.SetVersionTemplate(f.Version.Template())
	cmdutil.AddGroup(cmd, "Configuration & Setup",
		cmdconfig.NewCmdConfigure(f),
		cmdconfig.NewCmdUnconfigure(f),
		cmdconfig.NewCmdLogin(f),
		cmdconfig.NewCmdLogout(f),
	)
	cmdutil.AddGroup(cmd, "Devices",
		devices.NewCmdDevices(f),
		devices.NewCmdActivate(f),
		devices.NewCmdOpen(f),
		devices.NewCmdClose(f),
	)
	cmdutil.AddGroup(cmd, "Playback",
		player.NewCmdPlay(f),
		player.NewCmdPause(f),
		player.NewCmdNext(f),
		player.NewCmdPrevious(f),
		player.NewCmdStatus(f),
		volume.NewCmdVolume(f),
		queue.NewCmdQueue(f),
	)
	cmdutil.AddGroup(cmd, "Search",
		search.NewCmdSearch(f),
	)
	cmdutil.AddGroup(cmd, "Shelf",
		shelf.NewCmdShelf(f),
	)
	cmd.AddCommand(NewCmdVersion(f))
	cmd.AddCommand(NewCmdSplash(f))

	return cmd, nil
}

func runCmdRoot(cmd *cobra.Command) error {
	ui.ShowSplash()
	return cmd.Help()
}
