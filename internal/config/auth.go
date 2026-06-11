package config

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/matthew-collett/sp/internal/osutil"
	"golang.org/x/oauth2"
)

const (
	tokenFile   = "token"
	authURL     = "https://accounts.spotify.com/authorize"
	tokenURL    = "https://accounts.spotify.com/api/token"
	RedirectURL = "http://127.0.0.1:8080/callback"
)

var ErrNotAuthenticated = errors.New(`not authenticated. run "sp login"`)

var scopes = []string{
	"user-read-playback-state",
	"user-modify-playback-state",
	"user-read-currently-playing",
	"user-library-read",
	"user-library-modify",
	"user-read-recently-played",
	"playlist-read-private",
	"playlist-read-collaborative",
	"user-follow-read",
}

type Auth struct {
	tokenPath string
	oauth     *oauth2.Config
	source    oauth2.TokenSource
	token     *oauth2.Token
}

func NewAuth(clientID, clientSecret string) (*Auth, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	tokenPath := filepath.Join(dir, tokenFile)
	oauthCfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  RedirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	token := &oauth2.Token{}
	osutil.ReadFileSilent(tokenPath, token)
	ts := &tokenSource{
		source:    oauthCfg.TokenSource(context.Background(), token),
		tokenPath: tokenPath,
		last:      token.AccessToken,
	}
	return &Auth{
		tokenPath: tokenPath,
		oauth:     oauthCfg,
		source:    oauth2.ReuseTokenSource(token, ts),
		token:     token,
	}, nil
}

func (a *Auth) AuthURL(state string) string {
	return a.oauth.AuthCodeURL(state, oauth2.SetAuthURLParam("show_dialog", "true"))
}

func (a *Auth) Exchange(ctx context.Context, r *http.Request, state string) error {
	q := r.URL.Query()
	if e := q.Get("error"); e != "" {
		if e == "access_denied" {
			return errors.New("authentication cancelled by user")
		}
		return errors.New(e)
	}
	code := q.Get("code")
	if code == "" {
		return errors.New("no authorization code received")
	}
	if q.Get("state") != state {
		return errors.New("invalid state parameter")
	}
	token, err := a.oauth.Exchange(ctx, code)
	if err != nil {
		return err
	}
	a.token = token
	a.source = oauth2.ReuseTokenSource(token, &tokenSource{
		source:    a.oauth.TokenSource(ctx, token),
		tokenPath: a.tokenPath,
		last:      token.AccessToken,
	})
	return osutil.WriteFile(token, a.tokenPath, 0600)
}

func (a *Auth) HTTPClient(ctx context.Context) (*http.Client, error) {
	if !a.IsAuthenticated() {
		return nil, ErrNotAuthenticated
	}
	return oauth2.NewClient(ctx, a.source), nil
}

func (a *Auth) IsAuthenticated() bool {
	return a.token != nil && a.token.AccessToken != "" && a.token.RefreshToken != ""
}

func (a *Auth) Logout() error {
	a.token = &oauth2.Token{}
	if err := os.Remove(a.tokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

type tokenSource struct {
	source    oauth2.TokenSource
	tokenPath string
	last      string
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	token, err := s.source.Token()
	if err != nil {
		return nil, ErrNotAuthenticated
	}
	if token.AccessToken != s.last {
		s.last = token.AccessToken
		_ = osutil.WriteFile(token, s.tokenPath, 0600)
	}
	return token, nil
}
