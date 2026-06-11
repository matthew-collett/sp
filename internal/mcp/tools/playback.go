package tools

import (
	"context"
	"fmt"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type playback struct {
	f *factory.Factory
}

func RegisterPlayback(s *mcp.Server, f *factory.Factory) {
	t := &playback{f}
	mcp.AddTool(s, &mcp.Tool{Name: "get_status", Description: "Get the current Spotify playback status including track, artist, album, device, volume, progress, shuffle, and repeat state"}, t.status)
	mcp.AddTool(s, &mcp.Tool{Name: "play", Description: "Resume playback or play a shelf item by name. Pass name to play a specific shelf item, omit to resume current playback"}, t.play)
	mcp.AddTool(s, &mcp.Tool{Name: "pause", Description: "Pause Spotify playback"}, t.pause)
	mcp.AddTool(s, &mcp.Tool{Name: "next", Description: "Skip to the next track"}, t.next)
	mcp.AddTool(s, &mcp.Tool{Name: "previous", Description: "Go to the previous track"}, t.previous)
	mcp.AddTool(s, &mcp.Tool{Name: "seek", Description: "Seek to a position in the current track. Position format: 'm:ss' (e.g. '1:30') or seconds (e.g. '90')"}, t.seek)
}

func (t *playback) status(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	pb, err := sc.GetCurrentPlayback(ctx)
	if err != nil {
		return nil, nil, err
	}
	if pb.Item == nil {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No active playback"}}}, nil, nil
	}
	pct := 0
	if pb.Item.DurationMs > 0 {
		pct = int(float64(pb.ProgressMS) / float64(pb.Item.DurationMs) * 100)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(
		"Now playing: %s by %s\nAlbum: %s\nProgress: %s / %s (%d%%)\nDevice: %s (vol %d%%)\nShuffle: %v  Repeat: %s",
		pb.Item.Name,
		spotify.JoinArtists(pb.Item.Artists),
		pb.Item.Album.Name,
		spotify.FormatMS(pb.ProgressMS),
		spotify.FormatMS(pb.Item.DurationMs),
		pct,
		pb.Device.Name,
		pb.Device.VolumePercent,
		pb.ShuffleOn,
		pb.RepeatOn,
	)}}}, nil, nil
}

type playInput struct {
	Name string `json:"name"`
}

func (t *playback) play(ctx context.Context, _ *mcp.CallToolRequest, in playInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	var req spotify.PlayPlaybackRequest
	if in.Name != "" {
		shelf, err := t.f.Shelf()
		if err != nil {
			return nil, nil, err
		}
		item, ok := shelf.Get(in.Name)
		if !ok {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No shelf item named %q", in.Name)}}}, nil, nil
		}
		if spotify.IsTrack(item.URI) {
			req.URIs = []string{item.URI}
		} else {
			req.ContextURI = item.URI
		}
	}
	if err := sc.Play(ctx, device.ID, req); err != nil {
		return nil, nil, err
	}
	if in.Name != "" {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Playing %s", in.Name)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Playback resumed"}}}, nil, nil
}

func (t *playback) pause(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	if err := sc.Pause(ctx, device.ID); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Paused"}}}, nil, nil
}

func (t *playback) next(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	if err := sc.Next(ctx, device.ID); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Skipped to next track"}}}, nil, nil
}

func (t *playback) previous(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	if err := sc.Previous(ctx, device.ID); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Went to previous track"}}}, nil, nil
}

type seekInput struct {
	Position string `json:"position"`
}

func (t *playback) seek(ctx context.Context, _ *mcp.CallToolRequest, in seekInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	device, err := spotify.GetDevice(ctx, sc)
	if err != nil {
		return nil, nil, err
	}
	ms, err := spotify.ParsePosition(in.Position)
	if err != nil {
		return nil, nil, err
	}
	if err := sc.Seek(ctx, device.ID, ms); err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Seeked to %s", in.Position)}}}, nil, nil
}
