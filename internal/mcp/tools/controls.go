package tools

import (
	"context"
	"fmt"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type controls struct {
	f *factory.Factory
}

func RegisterControls(s *mcp.Server, f *factory.Factory) {
	t := &controls{f}
	mcp.AddTool(s, &mcp.Tool{Name: "toggle_shuffle", Description: "Toggle shuffle on or off for the current playback"}, t.shuffle)
	mcp.AddTool(s, &mcp.Tool{Name: "cycle_repeat", Description: "Cycle repeat mode: off → context (repeat playlist) → track (repeat one)"}, t.repeat)
	mcp.AddTool(s, &mcp.Tool{Name: "set_volume", Description: "Set the playback volume. volume is an integer from 0 to 100"}, t.volume)
	mcp.AddTool(s, &mcp.Tool{Name: "like", Description: "Save or unsave the currently playing track in the user's library. Toggles based on current saved state"}, t.like)
}

func (t *controls) shuffle(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	pb, err := sc.GetCurrentPlayback(ctx)
	if err != nil {
		return nil, nil, err
	}
	next := !pb.ShuffleOn
	if err := sc.SetShuffle(ctx, device.ID, next); err != nil {
		return nil, nil, err
	}
	if next {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Shuffle on"}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Shuffle off"}}}, nil, nil
}

func (t *controls) repeat(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	pb, err := sc.GetCurrentPlayback(ctx)
	if err != nil {
		return nil, nil, err
	}
	next := spotify.NextRepeatState(pb.RepeatOn)
	if err := sc.SetRepeat(ctx, device.ID, next); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Repeat %s", next)}}}, nil, nil
}

type volumeInput struct {
	Volume int `json:"volume"`
}

func (t *controls) volume(ctx context.Context, _ *mcp.CallToolRequest, in volumeInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	if err := sc.SetVolume(ctx, device.ID, in.Volume); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Volume set to %d%%", in.Volume)}}}, nil, nil
}

func (t *controls) like(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	pb, err := sc.GetCurrentPlayback(ctx)
	if err != nil {
		return nil, nil, err
	}
	if pb.Item == nil {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No track currently playing"}}}, nil, nil
	}
	saved, err := sc.IsTrackSaved(ctx, pb.Item.ID)
	if err != nil {
		return nil, nil, err
	}
	if saved {
		if err := sc.RemoveSavedTrack(ctx, pb.Item.ID); err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Unliked %s", pb.Item.Name)}}}, nil, nil
	}
	if err := sc.SaveTrack(ctx, pb.Item.ID); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Liked %s", pb.Item.Name)}}}, nil, nil
}
