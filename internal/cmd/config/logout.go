package config

import (
	"errors"

	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type logoutOpts struct {
	auth *config.Auth
}

func NewCmdLogout(f *factory.Factory) *cobra.Command {
	opts := &logoutOpts{}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out from Spotify",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			auth, err := f.Auth()
			if err != nil {
				return err
			}
			opts.auth = auth
			return runCmdLogout(opts)
		},
	}

	return cmd
}

func runCmdLogout(opts *logoutOpts) error {
	if err := opts.auth.Logout(); err != nil {
		return errors.New("failed to logout")
	}
	ui.Success("Bye!").Show()
	return nil
}
