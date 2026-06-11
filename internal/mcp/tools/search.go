package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type search struct {
	f *factory.Factory
}

func RegisterSearch(s *mcp.Server, f *factory.Factory) {
	t := &search{f}
	mcp.AddTool(s, &mcp.Tool{Name: "search", Description: "Search the Spotify catalog. type must be one of: track, album, playlist, artist. Returns up to 50 results"}, t.search)
}

type searchInput struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

func (t *search) search(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.Search(ctx, in.Query, []string{in.Type})
	if err != nil {
		return nil, nil, err
	}
	var sb strings.Builder
	switch in.Type {
	case "track":
		if res.Tracks == nil || len(res.Tracks.Items) == 0 {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No tracks found for %q", in.Query)}}}, nil, nil
		}
		fmt.Fprintf(&sb, "Tracks matching %q:\n", in.Query)
		for i, t := range res.Tracks.Items {
			fmt.Fprintf(&sb, "%d. %s by %s — %s\n", i+1, t.Name, spotify.JoinArtists(t.Artists), t.URI)
		}
	case "album":
		if res.Albums == nil || len(res.Albums.Items) == 0 {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No albums found for %q", in.Query)}}}, nil, nil
		}
		fmt.Fprintf(&sb, "Albums matching %q:\n", in.Query)
		for i, a := range res.Albums.Items {
			fmt.Fprintf(&sb, "%d. %s by %s — %s\n", i+1, a.Name, spotify.JoinArtists(a.Artists), a.URI)
		}
	case "playlist":
		if res.Playlists == nil || len(res.Playlists.Items) == 0 {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No playlists found for %q", in.Query)}}}, nil, nil
		}
		fmt.Fprintf(&sb, "Playlists matching %q:\n", in.Query)
		for i, p := range res.Playlists.Items {
			fmt.Fprintf(&sb, "%d. %s by %s — %s\n", i+1, p.Name, p.Owner.DisplayName, p.URI)
		}
	case "artist":
		if res.Artists == nil || len(res.Artists.Items) == 0 {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No artists found for %q", in.Query)}}}, nil, nil
		}
		fmt.Fprintf(&sb, "Artists matching %q:\n", in.Query)
		for i, a := range res.Artists.Items {
			fmt.Fprintf(&sb, "%d. %s — %s\n", i+1, a.Name, a.URI)
		}
	default:
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid type %q — must be one of: track, album, playlist, artist", in.Type)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}}}, nil, nil
}
