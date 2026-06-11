package spotify

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	reURI = regexp.MustCompile(`^spotify:(track|album|playlist|artist):[a-zA-Z0-9]+$`)

	ErrNoActiveDevice     = errors.New("no active playback device \nrun \"sp devices\" to see available devices\nrun \"sp activate <name>\" to activate one")
	ErrNoAvailableDevices = errors.New("no devices found\nrun \"sp open\" to open spotify on this device")
)

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

func ValidURI(uri string) bool {
	return reURI.MatchString(uri)
}

func IsTrack(uri string) bool {
	return strings.HasPrefix(uri, "spotify:track:")
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
