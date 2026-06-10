package config

import (
	"errors"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type unconfigureOpts struct {
	force bool
	cfg   *config.Config
}

func NewCmdUnconfigure(f *factory.Factory) *cobra.Command {
	opts := &unconfigureOpts{}

	cmd := &cobra.Command{
		Use:   "unconfigure",
		Short: "Remove Spotify API credentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.cfg = cfg
			return runCmdUnconfigure(opts)
		},
	}

	cmdutil.DisableConfigCheck(cmd)
	cmdutil.DisableAuthCheck(cmd)

	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "force removal without confirmation")

	return cmd
}

func runCmdUnconfigure(opts *unconfigureOpts) error {
	if !opts.cfg.IsConfigured() {
		ui.Warning("No configuration found").Show()
		return nil
	}

	if !opts.force {
		ui.Warning("This will remove your Spotify API credentials").Show()
		yes, err := ui.PromptConfirm("Are you sure?")
		if err != nil {
			return errors.New("prompt failed")
		}
		if !yes {
			ui.Info("Cancelled").Show()
			return nil
		}
	}

	if err := opts.cfg.Remove(); err != nil {
		return errors.New("failed to cleanup configuration")
	}

	ui.Success("Configuration removed").Show()
	return nil
}
