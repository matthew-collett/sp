package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/matthew-collett/sp/internal/build"
	"github.com/matthew-collett/sp/internal/cmd"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitAuth   exitCode = 2
	exitConfig exitCode = 3
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() exitCode {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cmdFactory := factory.New(build.Version)
	rootCmd, err := cmd.NewCmdRoot(cmdFactory)
	if err != nil {
		ui.Error(errors.New("failed to create root command")).Show()
		return exitError
	}

	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		ui.Error(err).Show()
		if errors.Is(err, config.ErrNoConfig) {
			return exitConfig
		}
		if errors.Is(err, config.ErrNotAuthenticated) {
			return exitAuth
		}
		return exitError
	}

	return exitOK
}
