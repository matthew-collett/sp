package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/matthew-collett/sp/internal/build"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/mcp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	f := factory.New(build.Version)
	return mcp.Serve(ctx, f)
}
