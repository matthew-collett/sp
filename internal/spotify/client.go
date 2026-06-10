package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"errors"
)

const (
	baseURL = "https://api.spotify.com/v1"
)

type Client struct {
	hc      *http.Client
	baseURL string
}

func NewClient(hc *http.Client) *Client {
	return &Client{hc: hc, baseURL: baseURL}
}

type RequestOptions struct {
	Headers map[string]string
	Params  url.Values
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var res User
	if err := c.get(ctx, me, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetDevices(ctx context.Context) (*DevicesResponse, error) {
	var res DevicesResponse
	if err := c.get(ctx, devices, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetPlaybackState(ctx context.Context) (*PlaybackState, error) {
	var res PlaybackState
	if err := c.get(ctx, player, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetCurrentPlayback(ctx context.Context) (*CurrentPlayback, error) {
	var res CurrentPlayback
	if err := c.get(ctx, player, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Play(ctx context.Context, deviceID string, req PlayPlaybackRequest) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	return c.put(ctx, play, req, opts, nil)
}

func (c *Client) Pause(ctx context.Context, deviceID string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	return c.put(ctx, pause, nil, opts, nil)
}

func (c *Client) Next(ctx context.Context, deviceID string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	return c.post(ctx, next, nil, opts, nil)
}

func (c *Client) Previous(ctx context.Context, deviceID string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	return c.post(ctx, previous, nil, opts, nil)
}

func (c *Client) SetVolume(ctx context.Context, deviceID string, volumePercent int) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	opts.Params.Set("volume_percent", strconv.Itoa(volumePercent))
	return c.put(ctx, volume, nil, opts, nil)
}

func (c *Client) SetShuffle(ctx context.Context, deviceID string, state bool) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	opts.Params.Set("state", strconv.FormatBool(state))
	return c.put(ctx, shuffle, nil, opts, nil)
}

func (c *Client) TransferPlayback(ctx context.Context, req TransferPlaybackRequest) error {
	return c.put(ctx, player, req, nil, nil)
}

func (c *Client) GetQueue(ctx context.Context) (*QueueResponse, error) {
	var res QueueResponse
	if err := c.get(ctx, queue, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) SetRepeat(ctx context.Context, deviceID string, state string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	opts.Params.Set("state", state)
	return c.put(ctx, repeat, nil, opts, nil)
}

func (c *Client) Seek(ctx context.Context, deviceID string, positionMS int) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("device_id", deviceID)
	opts.Params.Set("position_ms", strconv.Itoa(positionMS))
	return c.put(ctx, seek, nil, opts, nil)
}

func (c *Client) GetRecentlyPlayed(ctx context.Context) (*RecentlyPlayedResponse, error) {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("limit", "50")
	var res RecentlyPlayedResponse
	if err := c.get(ctx, recentlyPlayed, opts, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) SaveTrack(ctx context.Context, id string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("ids", id)
	return c.put(ctx, tracks, nil, opts, nil)
}

func (c *Client) RemoveSavedTrack(ctx context.Context, id string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("ids", id)
	return c.delete(ctx, tracks, opts, nil)
}

func (c *Client) IsTrackSaved(ctx context.Context, id string) (bool, error) {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("ids", id)
	var res []bool
	if err := c.get(ctx, tracks+"/contains", opts, &res); err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	return res[0], nil
}

func (c *Client) AddToQueue(ctx context.Context, uri string) error {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("uri", uri)
	return c.post(ctx, queue, nil, opts, nil)
}

func (c *Client) Search(ctx context.Context, query string, types []string) (*SearchResponse, error) {
	opts := &RequestOptions{Params: url.Values{}}
	opts.Params.Set("q", query)
	opts.Params.Set("type", strings.Join(types, ","))
	opts.Params.Set("limit", "50")

	var res SearchResponse
	if err := c.get(ctx, search, opts, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetSavedTracks(ctx context.Context) ([]SavedTrack, error) {
	return paginate(func(offset int) ([]SavedTrack, string, error) {
		opts := &RequestOptions{Params: url.Values{}}
		opts.Params.Set("offset", strconv.Itoa(offset))
		var res SavedTracksResponse
		if err := c.get(ctx, tracks, opts, &res); err != nil {
			return nil, "", err
		}
		return res.Items, res.Next, nil
	})
}

func (c *Client) GetSavedAlbums(ctx context.Context) ([]SavedAlbum, error) {
	return paginate(func(offset int) ([]SavedAlbum, string, error) {
		opts := &RequestOptions{Params: url.Values{}}
		opts.Params.Set("offset", strconv.Itoa(offset))
		var res SavedAlbumsResponse
		if err := c.get(ctx, albums, opts, &res); err != nil {
			return nil, "", err
		}
		return res.Items, res.Next, nil
	})
}

func (c *Client) GetPlaylists(ctx context.Context) ([]Playlist, error) {
	return paginate(func(offset int) ([]Playlist, string, error) {
		opts := &RequestOptions{Params: url.Values{}}
		opts.Params.Set("offset", strconv.Itoa(offset))
		var res PlaylistsResponse
		if err := c.get(ctx, playlists, opts, &res); err != nil {
			return nil, "", err
		}
		return res.Items, res.Next, nil
	})
}

func (c *Client) GetFollowedArtists(ctx context.Context) ([]Artist, error) {
	var all []Artist
	after := ""
	for {
		opts := &RequestOptions{Params: url.Values{}}
		opts.Params.Set("type", "artist")
		if after != "" {
			opts.Params.Set("after", after)
		}
		var res FollowedArtistsResponse
		if err := c.get(ctx, following, opts, &res); err != nil {
			return nil, err
		}
		all = append(all, res.Artists.Items...)
		if res.Artists.Next == "" {
			break
		}
		after = res.Artists.Cursors.After
	}
	return all, nil
}

func paginate[T any](fetch func(offset int) ([]T, string, error)) ([]T, error) {
	var all []T
	for offset := 0; ; {
		items, next, err := fetch(offset)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if next == "" {
			break
		}
		offset += len(items)
	}
	return all, nil
}

func (c *Client) get(ctx context.Context, path string, opts *RequestOptions, res any) error {
	resp, err := c.do(ctx, http.MethodGet, path, nil, opts)
	if err != nil {
		return err
	}
	return handleResponse(resp, res)
}

func (c *Client) post(ctx context.Context, path string, body any, opts *RequestOptions, res any) error {
	resp, err := c.do(ctx, http.MethodPost, path, body, opts)
	if err != nil {
		return err
	}
	return handleResponse(resp, res)
}

func (c *Client) put(ctx context.Context, path string, body any, opts *RequestOptions, res any) error {
	resp, err := c.do(ctx, http.MethodPut, path, body, opts)
	if err != nil {
		return err
	}
	return handleResponse(resp, res)
}

func (c *Client) delete(ctx context.Context, path string, opts *RequestOptions, res any) error {
	resp, err := c.do(ctx, http.MethodDelete, path, nil, opts)
	if err != nil {
		return err
	}
	return handleResponse(resp, res)
}

func (c *Client) do(ctx context.Context, method, path string, body any, opts *RequestOptions) (*http.Response, error) {
	reqURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	if opts != nil && opts.Params != nil {
		reqURL.RawQuery = opts.Params.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func handleResponse(resp *http.Response, res any) error {
	defer resp.Body.Close() //nolint:errcheck

	if isSpotifyError(resp) {
		return handleSpotifyError(resp)
	}

	if res != nil && resp.StatusCode != http.StatusNoContent {
		return json.NewDecoder(resp.Body).Decode(res)
	}
	return nil
}

func isSpotifyError(resp *http.Response) bool {
	return resp.StatusCode < 200 || resp.StatusCode >= 300
}

func handleSpotifyError(resp *http.Response) error {
	status := strings.ToLower(resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(status)
	}

	if len(body) > 0 {
		var sErr SpotifyError
		if err := json.Unmarshal(body, &sErr); err == nil {
			if strings.Contains(strings.ToLower(sErr.Error()), "restriction violated") {
				return nil
			}
			return sErr
		}
	}

	switch resp.StatusCode {
	case 405:
		if allow := resp.Header.Get("Allow"); allow != "" {
			return fmt.Errorf("%s (allowed: %s)", status, allow)
		}
	case 429:
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			return fmt.Errorf("%s (retry after %s seconds)", status, retryAfter)
		}
	}

	if len(body) > 0 {
		bodyStr := strings.TrimSpace(string(body))
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return fmt.Errorf("%s: %s", status, strings.ToLower(bodyStr))
	}

	return errors.New(status)
}
