package mcp

import (
	"context"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Serve(ctx context.Context, f *factory.Factory) error {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "sp",
		Version: f.Version.String(),
	}, nil)

	tools.RegisterPlayback(s, f)
	tools.RegisterControls(s, f)
	tools.RegisterQueue(s, f)
	tools.RegisterSearch(s, f)
	tools.RegisterDevices(s, f)
	tools.RegisterShelf(s, f)

	session, err := s.Connect(ctx, &mcp.StdioTransport{}, nil)
	if err != nil {
		return err
	}
	return session.Wait()
}
