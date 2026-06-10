package factory

import (
	"context"

	"github.com/matthew-collett/sp/internal/build"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/spotify"
)

type Factory struct {
	Version *build.VersionInfo
	Config  func() (*config.Config, error)
	Auth    func() (*config.Auth, error)
	Shelf   func() (*config.Shelf, error)
}

func New(version *build.VersionInfo) *Factory {
	f := &Factory{Version: version}
	f.Config = configFunc()
	f.Auth = authFunc(f)
	f.Shelf = shelfFunc()
	return f
}

func (f *Factory) SpotifyClient(ctx context.Context) (*spotify.Client, error) {
	auth, err := f.Auth()
	if err != nil {
		return nil, err
	}
	hc, err := auth.HTTPClient(ctx)
	if err != nil {
		return nil, err
	}
	return spotify.NewClient(hc), nil
}

func configFunc() func() (*config.Config, error) {
	var cached *config.Config
	var err error
	return func() (*config.Config, error) {
		if cached == nil && err == nil {
			cached, err = config.Load()
		}
		return cached, err
	}
}

func authFunc(f *Factory) func() (*config.Auth, error) {
	var cached *config.Auth
	var err error
	return func() (*config.Auth, error) {
		if cached == nil && err == nil {
			var cfg *config.Config
			if cfg, err = f.Config(); err == nil {
				cached, err = config.NewAuth(cfg.Credentials.ClientID, cfg.Credentials.ClientSecret)
			}
		}
		return cached, err
	}
}

func shelfFunc() func() (*config.Shelf, error) {
	var cached *config.Shelf
	var err error
	return func() (*config.Shelf, error) {
		if cached == nil && err == nil {
			cached, err = config.NewShelf()
		}
		return cached, err
	}
}
