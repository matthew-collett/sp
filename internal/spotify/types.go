package spotify

const (
	me                     = "/me"
	mePlayer               = "/me/player"
	mePlayerDevices        = "/me/player/devices"
	mePlayerPlay           = "/me/player/play"
	mePlayerPause          = "/me/player/pause"
	mePlayerNext           = "/me/player/next"
	mePlayerPrevious       = "/me/player/previous"
	mePlayerSeek           = "/me/player/seek"
	mePlayerVolume         = "/me/player/volume"
	mePlayerShuffle        = "/me/player/shuffle"
	mePlayerRepeat         = "/me/player/repeat"
	mePlayerQueue          = "/me/player/queue"
	mePlayerRecentlyPlayed = "/me/player/recently-played"
	meTracks               = "/me/tracks"
	meAlbums               = "/me/albums"
	mePlaylists            = "/me/playlists"
	meFollowing            = "/me/following"
	search                 = "/search"
	playlists              = "/playlists"
	artists                = "/artists"
)

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type TransferPlaybackRequest struct {
	DeviceIDs []string `json:"device_ids"`
	Play      bool     `json:"play,omitempty"`
}

type PlayPlaybackRequest struct {
	ContextURI string   `json:"context_uri,omitempty"`
	URIs       []string `json:"uris,omitempty"`
	Offset     *Offset  `json:"offset,omitempty"`
	PositionMS int      `json:"position_ms,omitempty"`
}

type Offset struct {
	Position int `json:"position,omitempty"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int    `json:"volume_percent"`
	SupportsVolume   bool   `json:"supports_volume"`
}

type VolumeRequest struct {
	VolumePercent int `json:"volume_percent"`
}

type PlaybackState struct {
	Device Device `json:"device"`
}

type CurrentPlayback struct {
	IsPlaying  bool   `json:"is_playing"`
	ProgressMS int    `json:"progress_ms"`
	ShuffleOn  bool   `json:"shuffle_state"`
	RepeatOn   string `json:"repeat_state"`
	Device     Device `json:"device"`
	Item       *Track `json:"item"`
}

type Track struct {
	ID         string   `json:"id"`
	URI        string   `json:"uri"`
	Name       string   `json:"name"`
	Artists    []Artist `json:"artists"`
	Album      Album    `json:"album"`
	DurationMs int      `json:"duration_ms"`
}

type Album struct {
	ID          string   `json:"id"`
	URI         string   `json:"uri"`
	Name        string   `json:"name"`
	Artists     []Artist `json:"artists"`
	TotalTracks int      `json:"total_tracks"`
	Images      []Image  `json:"images"`
}

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PlaylistTracks struct {
	Total int `json:"total"`
}

type Playlist struct {
	ID     string         `json:"id"`
	URI    string         `json:"uri"`
	Name   string         `json:"name"`
	Public bool           `json:"public"`
	Owner  Owner          `json:"owner"`
	Tracks PlaylistTracks `json:"tracks"`
}

type Artist struct {
	ID        string    `json:"id"`
	URI       string    `json:"uri"`
	Name      string    `json:"name"`
	Genres    []string  `json:"genres"`
	Followers Followers `json:"followers"`
}

type Followers struct {
	Total int `json:"total"`
}

type Owner struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type SearchResponse struct {
	Tracks    *TracksResponse    `json:"tracks,omitempty"`
	Albums    *AlbumsResponse    `json:"albums,omitempty"`
	Playlists *PlaylistsResponse `json:"playlists,omitempty"`
	Artists   *ArtistsResponse   `json:"artists,omitempty"`
	Total     int                `json:"total,omitempty"`
}

type TracksResponse struct {
	Items []Track `json:"items"`
	Total int     `json:"total"`
}

type AlbumsResponse struct {
	Items []Album `json:"items"`
	Total int     `json:"total"`
}

type PlaylistsResponse struct {
	Items []Playlist `json:"items"`
	Total int        `json:"total"`
	Next  string     `json:"next"`
}

type ArtistsResponse struct {
	Items []Artist `json:"items"`
	Total int      `json:"total"`
}

type SavedTrack struct {
	AddedAt string `json:"added_at"`
	Track   Track  `json:"track"`
}

type SavedAlbum struct {
	AddedAt string `json:"added_at"`
	Album   Album  `json:"album"`
}

type SavedTracksResponse struct {
	Items []SavedTrack `json:"items"`
	Total int          `json:"total"`
	Next  string       `json:"next"`
}

type SavedAlbumsResponse struct {
	Items []SavedAlbum `json:"items"`
	Total int          `json:"total"`
	Next  string       `json:"next"`
}

type ArtistCursors struct {
	After string `json:"after"`
}

type ArtistPage struct {
	Items   []Artist      `json:"items"`
	Total   int           `json:"total"`
	Next    string        `json:"next"`
	Cursors ArtistCursors `json:"cursors"`
}

type FollowedArtistsResponse struct {
	Artists ArtistPage `json:"artists"`
}

type QueueResponse struct {
	CurrentlyPlaying *Track  `json:"currently_playing"`
	Queue            []Track `json:"queue"`
}

type PlayHistoryItem struct {
	Track    Track  `json:"track"`
	PlayedAt string `json:"played_at"`
}

type RecentlyPlayedResponse struct {
	Items []PlayHistoryItem `json:"items"`
}
