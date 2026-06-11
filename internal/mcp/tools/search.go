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
	mcp.AddTool(s, &mcp.Tool{Name: "search", Description: "Search Spotify catalog or your library. type must be one of: track, album, playlist, artist. Set mine=true to search only your saved library"}, t.search)
}

type searchInput struct {
	Query string `json:"query"`
	Type  string `json:"type"`
	Mine  bool   `json:"mine"`
}

func (t *search) search(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, any, error) {
	sc, err := t.f.SpotifyClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	if in.Mine {
		return searchLibrary(ctx, sc, in)
	}
	return searchCatalog(ctx, sc, in)
}

func searchCatalog(ctx context.Context, sc *spotify.Client, in searchInput) (*mcp.CallToolResult, any, error) {
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

func searchLibrary(ctx context.Context, sc *spotify.Client, in searchInput) (*mcp.CallToolResult, any, error) {
	q := strings.ToLower(in.Query)
	var sb strings.Builder
	n := 0
	switch in.Type {
	case "track":
		items, err := sc.GetSavedTracks(ctx)
		if err != nil {
			return nil, nil, err
		}
		fmt.Fprintf(&sb, "Saved tracks:\n")
		for _, s := range items {
			if q != "" && !strings.Contains(strings.ToLower(s.Track.Name), q) {
				continue
			}
			n++
			fmt.Fprintf(&sb, "%d. %s by %s — %s\n", n, s.Track.Name, spotify.JoinArtists(s.Track.Artists), s.Track.URI)
		}
	case "album":
		items, err := sc.GetSavedAlbums(ctx)
		if err != nil {
			return nil, nil, err
		}
		fmt.Fprintf(&sb, "Saved albums:\n")
		for _, s := range items {
			if q != "" && !strings.Contains(strings.ToLower(s.Album.Name), q) {
				continue
			}
			n++
			fmt.Fprintf(&sb, "%d. %s by %s — %s\n", n, s.Album.Name, spotify.JoinArtists(s.Album.Artists), s.Album.URI)
		}
	case "playlist":
		items, err := sc.GetPlaylists(ctx)
		if err != nil {
			return nil, nil, err
		}
		fmt.Fprintf(&sb, "Your playlists:\n")
		for _, p := range items {
			if q != "" && !strings.Contains(strings.ToLower(p.Name), q) {
				continue
			}
			n++
			fmt.Fprintf(&sb, "%d. %s by %s — %s\n", n, p.Name, p.Owner.DisplayName, p.URI)
		}
	case "artist":
		items, err := sc.GetFollowedArtists(ctx)
		if err != nil {
			return nil, nil, err
		}
		fmt.Fprintf(&sb, "Followed artists:\n")
		for _, a := range items {
			if q != "" && !strings.Contains(strings.ToLower(a.Name), q) {
				continue
			}
			n++
			fmt.Fprintf(&sb, "%d. %s — %s\n", n, a.Name, a.URI)
		}
	default:
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid type %q — must be one of: track, album, playlist, artist", in.Type)}}}, nil, nil
	}
	if n == 0 {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No %ss found in your library", in.Type)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: sb.String()}}}, nil, nil
}
