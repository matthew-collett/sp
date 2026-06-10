package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/matthew-collett/sp/internal/osutil"
)

const (
	dirName    = ".sp"
	configFile = "config"
)

var ErrNoConfig = errors.New(`spotify credentials not configured. run "sp configure"`)

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Config struct {
	path        string
	Credentials Credentials `json:"credentials"`
	User        string      `json:"user,omitempty"`
}

func Load() (*Config, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	c := &Config{path: filepath.Join(dir, configFile)}
	osutil.ReadFileSilent(c.path, c)
	return c, nil
}

func (c *Config) Write() error {
	return osutil.WriteFile(c, c.path, 0644)
}

func (c *Config) Remove() error {
	c.Credentials = Credentials{}
	if err := os.Remove(c.path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c *Config) IsConfigured() bool {
	return c.Credentials.ClientID != "" && c.Credentials.ClientSecret != ""
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, dirName), nil
}
