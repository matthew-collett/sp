package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type queue struct {
	f *factory.Factory
}

func RegisterQueue(s *mcp.Server, f *factory.Factory) {
	t := &queue{f}
	mcp.AddTool(s, &mcp.Tool{Name: "get_queue", Description: "Get the current Spotify playback queue including the currently playing track and upcoming tracks"}, t.getQueue)
	mcp.AddTool(s, &mcp.Tool{Name: "add_to_queue", Description: "Add a track to the playback queue by Spotify URI or shelf name — only tracks can be queued"}, t.addToQueue)
	mcp.AddTool(s, &mcp.Tool{Name: "get_recent", Description: "Get the 50 most recently played tracks"}, t.getRecent)
}

func (t *queue) getQueue(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.GetQueue(ctx)
	if err != nil {
		return nil, nil, err
	}
	var sb strings.Builder
	if res.CurrentlyPlaying != nil {
		fmt.Fprintf(&sb, "Now playing: %s by %s\n\n", res.CurrentlyPlaying.Name, spotify.JoinArtists(res.CurrentlyPlaying.Artists))
	}
	if len(res.Queue) == 0 {
		sb.WriteString("Queue is empty")
	} else {
		sb.WriteString("Up next:\n")
		for i, track := range res.Queue {
			fmt.Fprintf(&sb, "%d. %s by %s\n", i+1, track.Name, spotify.JoinArtists(track.Artists))
		}
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}}}, nil, nil
}

type addToQueueInput struct {
	URI string `json:"uri"`
}

func (t *queue) addToQueue(ctx context.Context, _ *mcp.CallToolRequest, in addToQueueInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	uri := in.URI
	var name string
	if _, _, err := spotify.ParseURI(in.URI); err != nil {
		shelf, err := t.f.Shelf()
		if err != nil {
			return nil, nil, err
		}
		item, ok := shelf.Get(in.URI)
		if !ok {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No shelf item named %q", in.URI)}}}, nil, nil
		}
		uri = item.URI
		name = item.Name
	}
	if kind, _, _ := spotify.ParseURI(uri); kind != "track" {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Only tracks can be added to the queue"}}}, nil, nil
	}
	if err := sc.AddToQueue(ctx, uri); err != nil {
		return nil, nil, err
	}
	if name == "" {
		if q, err := sc.GetQueue(ctx); err == nil {
			for _, t := range q.Queue {
				if t.URI == uri {
					name = t.Name
					break
				}
			}
		}
	}
	if name != "" {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Added %q to queue", name)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Added to queue"}}}, nil, nil
}

func (t *queue) getRecent(ctx context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.GetRecentlyPlayed(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(res.Items) == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "No recently played tracks"}}}, nil, nil
	}
	var sb strings.Builder
	sb.WriteString("Recently played:\n")
	for i, item := range res.Items {
		fmt.Fprintf(&sb, "%d. %s by %s\n", i+1, item.Track.Name, spotify.JoinArtists(item.Track.Artists))
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}}}, nil, nil
}
