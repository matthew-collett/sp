package config

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "embed"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/osutil"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

//go:embed templates/auth-success.html
var authSuccessHTML string

//go:embed templates/auth-error.html
var authErrorHTML string

type loginOpts struct {
	noBrowser bool
	force     bool
	auth      *config.Auth
	cfg       *config.Config
}

func NewCmdLogin(f *factory.Factory) *cobra.Command {
	opts := &loginOpts{}

	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.NoArgs,
		Short: "Authenticate with your Spotify account",
		RunE: func(cmd *cobra.Command, args []string) error {
			auth, err := f.Auth()
			if err != nil {
				return err
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			opts.auth = auth
			opts.cfg = cfg
			return runCmdLogin(cmd.Context(), f, opts)
		},
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "force authentication")
	cmd.Flags().BoolVarP(&opts.noBrowser, "no-browser", "n", false, "don't open browser automatically")

	return cmd
}

func runCmdLogin(ctx context.Context, f *factory.Factory, opts *loginOpts) error {
	if opts.auth.IsAuthenticated() && !opts.force {
		ui.Success("Already authenticated with Spotify!").Show()
		yes, err := ui.PromptConfirm("Would you like to re-authenticate?")
		if err != nil {
			return errors.New("prompt failed")
		}
		if !yes {
			ui.Success("Using existing authentication").Show()
			return nil
		}
	}

	state, err := generateState()
	if err != nil {
		return errors.New("failed to start login flow")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	tokenChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		err := opts.auth.Exchange(ctx, r, state)
		w.Header().Set("Content-Type", "text/html")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, authErrorHTML) //nolint:errcheck
			tokenChan <- err
			return
		}
		fmt.Fprint(w, authSuccessHTML) //nolint:errcheck
		tokenChan <- nil
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go server.ListenAndServe() //nolint:errcheck
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx) //nolint:errcheck
	}()

	time.Sleep(100 * time.Millisecond)

	authURL := opts.auth.AuthURL(state)
	if !opts.noBrowser {
		ui.Info("Opening browser for Spotify authentication...").Show()
		if err := osutil.OpenBrowser(authURL); err != nil {
			ui.Warning("Failed to open browser automatically").Show()
			ui.Info("Please visit: %s", ui.Text(authURL).Underline()).Show()
		}
	} else {
		ui.Info("Visit this URL to authenticate: %s", ui.Text(authURL).Underline()).Show()
	}

	select {
	case err := <-tokenChan:
		if err != nil {
			return errors.New("authentication cancelled")
		}
		if sc, err := f.SpotifyClient(ctx); err == nil {
			if resp, err := sc.GetCurrentUser(ctx); err == nil {
				user := resp.DisplayName
				if user == "" {
					user = resp.Email
				}
				if user == "" {
					user = resp.ID
				}
				opts.cfg.User = user
				if err := opts.cfg.Write(); err != nil {
					return err
				}
			}
		}
		ui.Success("Successfully authenticated with Spotify!").Show()
		return nil
	case <-ctx.Done():
		return errors.New("authentication timed out")
	}
}

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
