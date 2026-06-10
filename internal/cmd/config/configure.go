package config

import (
	"errors"
	"strings"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

const (
	minCredentialLength = 10
	spotifyDashboardURL = "https://developer.spotify.com/dashboard"
)

type configureOpts struct {
	force        bool
	clientID     string
	clientSecret string
	cfg          *config.Config
}

func NewCmdConfigure(f *factory.Factory) *cobra.Command {
	opts := &configureOpts{}

	cmd := &cobra.Command{
		Use:   "configure",
		Args:  cobra.NoArgs,
		Short: "Configure sp with Spotify API credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.cfg = cfg
			return runCmdConfigure(opts)
		},
	}

	cmdutil.DisableConfigCheck(cmd)
	cmdutil.DisableAuthCheck(cmd)

	cmd.Flags().StringVarP(&opts.clientID, "client-id", "c", "", "spotify client ID")
	cmd.Flags().StringVarP(&opts.clientSecret, "client-secret", "s", "", "spotify client secret")
	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "force overwrite existing configuration")

	return cmd
}

func runCmdConfigure(opts *configureOpts) error {
	if !opts.force && opts.cfg.IsConfigured() {
		ui.Success("sp is already configured!").Show()
		ui.Info("Client ID: %s", opts.cfg.Credentials.ClientID).Show()
		ui.Info("Client Secret: %s", strings.Repeat("*", len(opts.cfg.Credentials.ClientSecret))).Show()

		yes, err := ui.PromptConfirm("Would you like to reconfigure?")
		if err != nil {
			return err
		}
		if !yes {
			ui.Success("Using existing configuration").Show()
			return nil
		}
	}
	ui.Text("sp Configuration").Bold().Show()
	ui.Info("You need to register your app on Spotify Developer Dashboard:").Show()
	ui.Info("1. Go to: %s", ui.Text(spotifyDashboardURL).Underline()).Tab().Show()
	ui.Info("2. Create a new app").Tab().Show()
	ui.Info("3. Set the redirect URI to: %s", ui.Text(config.RedirectURL).Underline()).Tab().Show()

	if opts.clientID == "" {
		clientID, err := ui.Prompt("Enter your Spotify Client ID")
		if err != nil {
			return errors.New("prompt failed")
		}
		opts.clientID = clientID
	}
	if opts.clientID == "" {
		return errors.New("client ID cannot be empty")
	}

	if opts.clientSecret == "" {
		clientSecret, err := ui.PromptSecret("Enter your Spotify Client Secret")
		if err != nil {
			return errors.New("prompt failed")
		}
		opts.clientSecret = clientSecret
	}
	if opts.clientSecret == "" {
		return errors.New("client secret cannot be empty")
	}

	if len(opts.clientID) < minCredentialLength || len(opts.clientSecret) < minCredentialLength {
		return errors.New("credentials seem too short. Please check and try again")
	}

	opts.cfg.Credentials = config.Credentials{
		ClientID:     opts.clientID,
		ClientSecret: opts.clientSecret,
	}

	if err := opts.cfg.Write(); err != nil {
		return errors.New("failed to save configuration")
	}

	dir, err := config.Dir()
	if err != nil {
		return err
	}
	ui.Success("Configuration saved to %s", dir).Show()
	ui.Info(`Run "sp login" to authenticate with Spotify.`).Show()
	return nil
}
