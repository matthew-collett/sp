package config

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/matthew-collett/sp/internal/osutil"
	"golang.org/x/oauth2"
)

const (
	tokenFile   = "token"
	authURL     = "https://accounts.spotify.com/authorize"
	tokenURL    = "https://accounts.spotify.com/api/token"
	RedirectURL = "http://127.0.0.1:8080/callback"
)

var (
	ErrNotAuthenticated = errors.New(`not authenticated. run "sp login"`)
	ErrInvalidState     = errors.New("invalid state parameter")
)

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
	source    *tokenSource
	verifier  string
}

func NewAuth(clientID string) (*Auth, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	tokenPath := filepath.Join(dir, tokenFile)
	oauthCfg := &oauth2.Config{
		ClientID:    clientID,
		RedirectURL: RedirectURL,
		Scopes:      scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	return &Auth{
		tokenPath: tokenPath,
		oauth:     oauthCfg,
		source:    &tokenSource{oauth: oauthCfg, tokenPath: tokenPath},
	}, nil
}

func (a *Auth) AuthURL(state string) string {
	a.verifier = oauth2.GenerateVerifier()
	return a.oauth.AuthCodeURL(state,
		oauth2.S256ChallengeOption(a.verifier),
		oauth2.SetAuthURLParam("show_dialog", "true"),
	)
}

func (a *Auth) Exchange(ctx context.Context, r *http.Request, state string) error {
	q := r.URL.Query()
	if q.Get("state") != state {
		return ErrInvalidState
	}
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
	token, err := a.oauth.Exchange(ctx, code, oauth2.VerifierOption(a.verifier))
	if err != nil {
		return err
	}
	return osutil.WriteFile(token, a.tokenPath, 0600)
}

func (a *Auth) HTTPClient(ctx context.Context) (*http.Client, error) {
	if !a.IsAuthenticated() {
		return nil, ErrNotAuthenticated
	}
	return oauth2.NewClient(ctx, a.source), nil
}

func (a *Auth) IsAuthenticated() bool {
	token := readToken(a.tokenPath)
	return token.AccessToken != "" && token.RefreshToken != ""
}

func (a *Auth) Logout() error {
	if err := os.Remove(a.tokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

type tokenSource struct {
	mu        sync.Mutex
	oauth     *oauth2.Config
	tokenPath string
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := readToken(s.tokenPath)
	if token.RefreshToken == "" {
		return nil, ErrNotAuthenticated
	}
	if token.Valid() {
		return token, nil
	}

	refreshed, err := s.oauth.TokenSource(context.Background(), token).Token()
	if err != nil {
		if current := readToken(s.tokenPath); current.Valid() {
			return current, nil
		}
		return nil, ErrNotAuthenticated
	}
	if err := osutil.OverwriteFile(refreshed, s.tokenPath); err != nil {
		return nil, ErrNotAuthenticated
	}
	return refreshed, nil
}

func readToken(path string) *oauth2.Token {
	token := &oauth2.Token{}
	osutil.ReadFileSilent(path, token)
	return token
}
