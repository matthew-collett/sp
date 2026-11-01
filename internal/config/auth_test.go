package config

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matthew-collett/sp/internal/osutil"
	"golang.org/x/oauth2"
)

func newTestAuth(t *testing.T, handler http.HandlerFunc) *Auth {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	tokenPath := filepath.Join(t.TempDir(), tokenFile)
	oauthCfg := &oauth2.Config{
		ClientID:    "client",
		RedirectURL: RedirectURL,
		Scopes:      scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   server.URL + "/authorize",
			TokenURL:  server.URL + "/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	return &Auth{
		tokenPath: tokenPath,
		oauth:     oauthCfg,
		source:    &tokenSource{oauth: oauthCfg, tokenPath: tokenPath},
	}
}

func tokenHandler(access, refresh string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("client_secret") != "" {
			http.Error(w, `{"error":"invalid_client"}`, http.StatusUnauthorized)
			return
		}
		if r.FormValue("grant_type") == "authorization_code" && r.FormValue("code_verifier") == "" {
			http.Error(w, `{"error":"invalid_grant"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"` + access + `","refresh_token":"` + refresh + `","token_type":"Bearer","expires_in":3600}`)) //nolint:errcheck
	}
}

func callback(q url.Values) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/callback?"+q.Encode(), nil)
}

func writeTestToken(t *testing.T, path string, expiry time.Time) {
	t.Helper()
	token := &oauth2.Token{AccessToken: "access", RefreshToken: "refresh", TokenType: "Bearer", Expiry: expiry}
	if err := osutil.WriteFile(token, path, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestAuthURLUsesPKCE(t *testing.T) {
	auth := newTestAuth(t, tokenHandler("access", "refresh"))

	u, err := url.Parse(auth.AuthURL("state"))
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()

	if q.Get("code_challenge_method") != "S256" {
		t.Error("auth URL should use the S256 code challenge method")
	}
	if q.Get("code_challenge") == "" {
		t.Error("auth URL should include a code challenge")
	}
	if q.Get("redirect_uri") != RedirectURL {
		t.Error("auth URL should use the loopback redirect URL")
	}
}

func TestExchangeRejectsWrongState(t *testing.T) {
	auth := newTestAuth(t, tokenHandler("access", "refresh"))
	auth.AuthURL("state")

	err := auth.Exchange(context.Background(), callback(url.Values{"state": {"other"}, "code": {"c"}}), "state")
	if !errors.Is(err, ErrInvalidState) {
		t.Errorf("got %v, want ErrInvalidState", err)
	}
	if auth.IsAuthenticated() {
		t.Error("wrong state should not authenticate")
	}
}

func TestExchangeChecksStateBeforeError(t *testing.T) {
	auth := newTestAuth(t, tokenHandler("access", "refresh"))
	auth.AuthURL("state")

	err := auth.Exchange(context.Background(), callback(url.Values{"state": {"other"}, "error": {"access_denied"}}), "state")
	if !errors.Is(err, ErrInvalidState) {
		t.Errorf("got %v, want ErrInvalidState", err)
	}
}

func TestExchangeStoresToken(t *testing.T) {
	auth := newTestAuth(t, tokenHandler("access", "refresh"))
	auth.AuthURL("state")

	if err := auth.Exchange(context.Background(), callback(url.Values{"state": {"state"}, "code": {"c"}}), "state"); err != nil {
		t.Fatal(err)
	}
	if !auth.IsAuthenticated() {
		t.Error("valid callback should authenticate")
	}
	info, err := os.Stat(auth.tokenPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("token file mode is %o, want 600", info.Mode().Perm())
	}
}

func TestTokenSourceRefreshesExpiredToken(t *testing.T) {
	auth := newTestAuth(t, tokenHandler("refreshed", "rotated"))
	writeTestToken(t, auth.tokenPath, time.Now().Add(-time.Hour))

	token, err := auth.source.Token()
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "refreshed" {
		t.Errorf("got access token %q, want refreshed", token.AccessToken)
	}
	if stored := readToken(auth.tokenPath); stored.RefreshToken != "rotated" {
		t.Errorf("stored refresh token is %q, want rotated", stored.RefreshToken)
	}
}

func TestTokenSourceHonoursLogout(t *testing.T) {
	auth := newTestAuth(t, tokenHandler("access", "refresh"))
	writeTestToken(t, auth.tokenPath, time.Now().Add(time.Hour))

	if err := auth.Logout(); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.source.Token(); !errors.Is(err, ErrNotAuthenticated) {
		t.Errorf("got %v, want ErrNotAuthenticated", err)
	}
}

func TestTokenSourceDoesNotRecreateTokenAfterLogout(t *testing.T) {
	var auth *Auth
	refresh := tokenHandler("refreshed", "rotated")
	auth = newTestAuth(t, func(w http.ResponseWriter, r *http.Request) {
		auth.Logout() //nolint:errcheck
		refresh(w, r)
	})
	writeTestToken(t, auth.tokenPath, time.Now().Add(-time.Hour))

	if _, err := auth.source.Token(); !errors.Is(err, ErrNotAuthenticated) {
		t.Errorf("got %v, want ErrNotAuthenticated", err)
	}
	if _, err := os.Stat(auth.tokenPath); !os.IsNotExist(err) {
		t.Error("refresh should not recreate a logged out token")
	}
}

func TestConfigRemoveDeletesToken(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{path: filepath.Join(dir, configFile), Credentials: Credentials{ClientID: "client"}}
	if err := cfg.Write(); err != nil {
		t.Fatal(err)
	}
	writeTestToken(t, filepath.Join(dir, tokenFile), time.Now().Add(time.Hour))

	if err := cfg.Remove(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, tokenFile)); !os.IsNotExist(err) {
		t.Error("remove should delete the token")
	}
}
