package spotify

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	reURI = regexp.MustCompile(`^spotify:(track|album|playlist|artist):([a-zA-Z0-9]+)$`)

	ErrNoActiveDevice     = errors.New("no active playback device \nrun \"sp devices\" to see available devices\nrun \"sp activate <name>\" to activate one")
	ErrNoAvailableDevices = errors.New("no devices found\nrun \"sp open\" to open spotify on this device")
)

const playbackDelay = 500 * time.Millisecond

func GetDevice(ctx context.Context, sc *Client) (*Device, error) {
	resp, err := sc.GetDevices(ctx)
	if err != nil {
		return nil, err
	}
	if len(resp.Devices) == 0 {
		return nil, ErrNoAvailableDevices
	}
	for i := range resp.Devices {
		if resp.Devices[i].IsActive {
			return &resp.Devices[i], nil
		}
	}
	return nil, ErrNoActiveDevice
}

func ParseURI(uri string) (kind, id string, err error) {
	if m := reURI.FindStringSubmatch(uri); m != nil {
		return m[1], m[2], nil
	}
	return "", "", fmt.Errorf("invalid spotify uri: %q", uri)
}

func GetCurrentPlaybackDelayed(ctx context.Context, sc *Client) (*CurrentPlayback, error) {
	time.Sleep(playbackDelay)
	return sc.GetCurrentPlayback(ctx)
}

func URIName(ctx context.Context, sc *Client, uri string) string {
	kind, id, err := ParseURI(uri)
	if err != nil {
		return ""
	}
	if kind == "playlist" {
		if pl, err := sc.GetPlaylist(ctx, id); err == nil {
			return pl.Name
		}
		return ""
	}
	pb, err := GetCurrentPlaybackDelayed(ctx, sc)
	if err != nil || pb.Item == nil {
		return ""
	}
	if kind == "album" {
		return pb.Item.Album.Name
	}
	return pb.Item.Name
}

func JoinArtists(artists []Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

func FormatMS(ms int) string {
	s := ms / 1000
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func ParsePosition(s string) (int, error) {
	if strings.Contains(s, ":") {
		var mins, secs int
		if _, err := fmt.Sscanf(s, "%d:%d", &mins, &secs); err != nil {
			return 0, fmt.Errorf("invalid position %q", s)
		}
		return (mins*60 + secs) * 1000, nil
	}
	var secs int
	if _, err := fmt.Sscanf(s, "%d", &secs); err != nil {
		return 0, fmt.Errorf("invalid position %q — use m:ss or seconds", s)
	}
	return secs * 1000, nil
}

func NextRepeatState(current string) string {
	switch current {
	case "off":
		return "context"
	case "context":
		return "track"
	default:
		return "off"
	}
}
