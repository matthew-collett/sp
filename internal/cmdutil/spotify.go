package cmdutil

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/matthew-collett/sp/internal/spotify"
)

var (
	reURI = regexp.MustCompile(`^spotify:(track|album|playlist|artist):[a-zA-Z0-9]+$`)

	ErrNoActiveDevice     = errors.New("no active playback device \nrun \"sp devices\" to see available devices\nrun \"sp activate <name>\" to activate one")
	ErrNoAvailableDevices = errors.New("no devices found\nrun \"sp open\" to open spotify on this device")
)

func GetDevice(ctx context.Context, sc *spotify.Client) (*spotify.Device, error) {
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

func JoinArtists(artists []spotify.Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}
