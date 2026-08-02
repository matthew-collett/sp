package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type devices struct {
	f *factory.Factory
}

func RegisterDevices(s *mcp.Server, f *factory.Factory) {
	t := &devices{f}
	mcp.AddTool(s, &mcp.Tool{Name: "get_devices", Description: "List all available Spotify playback devices"}, t.getDevices)
	mcp.AddTool(s, &mcp.Tool{Name: "activate_device", Description: "Switch Spotify playback to a device by name (case-insensitive partial match)"}, t.activateDevice)
}

func (t *devices) getDevices(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.GetDevices(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(res.Devices) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No devices available"}}}, nil, nil
	}
	var sb strings.Builder
	sb.WriteString("Available devices:\n")
	for _, d := range res.Devices {
		active := ""
		if d.IsActive {
			active = " (active)"
		}
		fmt.Fprintf(&sb, "- %s [%s]%s\n", d.Name, d.Type, active)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}}}, nil, nil
}

type activateInput struct {
	Name string `json:"name"`
}

func (t *devices) activateDevice(ctx context.Context, _ *mcp.CallToolRequest, in activateInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.GetDevices(ctx)
	if err != nil {
		return nil, nil, err
	}
	query := strings.ToLower(in.Name)
	for _, d := range res.Devices {
		if strings.Contains(strings.ToLower(d.Name), query) {
			if err := sc.TransferPlayback(ctx, spotify.TransferPlaybackRequest{DeviceIDs: []string{d.ID}, Play: true}); err != nil {
				return nil, nil, err
			}
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Switched to %q", d.Name)}}}, nil, nil
		}
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No device found matching %q", in.Name)}}}, nil, nil
}
