package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type shelf struct {
	f *factory.Factory
}

func RegisterShelf(s *mcp.Server, f *factory.Factory) {
	t := &shelf{f}
	mcp.AddTool(s, &mcp.Tool{Name: "list_shelf", Description: "List all shelf shortcuts — named Spotify URIs the user has saved for quick playback"}, t.list)
}

func (t *shelf) list(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sh, err := t.f.Shelf()
	if err != nil {
		return nil, nil, err
	}
	if len(sh.Items) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Shelf is empty"}}}, nil, nil
	}
	var sb strings.Builder
	sb.WriteString("Shelf items:\n")
	for name, item := range sh.Items {
		fmt.Fprintf(&sb, "- %s: %s\n", name, item.URI)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}}}, nil, nil
}
