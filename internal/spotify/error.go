package spotify

import (
	"fmt"
	"strings"
)

type SpotifyError struct {
	Err `json:"error"`
}

type Err struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
}

func (e SpotifyError) Error() string {
	return fmt.Sprintf("Spotify error: %s", strings.ToLower(e.Message))
}

func (e SpotifyError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprint(s, formatErr(e)) //nolint:errcheck
			return
		}
		fallthrough
	case 's', 'q':
		fmt.Fprint(s, e.Error()) //nolint:errcheck
	}
}

func formatErr(err SpotifyError) string {
	if err.Reason != "" {
		return fmt.Sprintf("HTTP %d, message: %s, reason: %s", err.Status, err.Message, err.Reason)
	}
	return fmt.Sprintf("HTTP %d: %s", err.Status, err.Message)
}
