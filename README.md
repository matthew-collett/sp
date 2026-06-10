<div align="center">

<img width="400" src="etc/img.svg" alt="sp demo" />

**A fast, minimal Spotify CLI for your terminal**

[![Release](https://img.shields.io/github/v/release/matthew-collett/sp)](https://github.com/matthew-collett/sp/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/matthew-collett/sp)](go.mod)
[![License](https://img.shields.io/github/license/matthew-collett/sp)](LICENSE)

Control playback, search your library, manage devices, and save quick-play shortcuts — all without leaving the terminal.

<img src="etc/demo.gif" alt="sp demo" />

</div>

## Installation

```sh
brew tap matthew-collett/sp
brew trust matthew-collett/sp
brew install sp
```

<details>
<summary>Install from source</summary>

```sh
git clone https://github.com/matthew-collett/sp
cd sp
make install
```

</details>

## Setup

`sp` uses the Spotify Web API and requires a developer application.

1. Create an app at the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Add `http://localhost:8080/callback` as a redirect URI
3. Copy your Client ID and Client Secret

```sh
sp configure   # enter your credentials
sp login       # authenticate via browser
```

## Useful Commands

### Playback

```sh
sp status                  # show now playing, progress, shuffle, repeat
sp play lofi               # play a shelf item by name
sp play lofi --shuffle     # play shuffled
sp like                    # save or unsave the current track
sp shuffle                 # toggle shuffle
sp repeat                  # cycle repeat (off → context → track)
sp seek 1:30               # seek to 1m30s in the current track
sp recent                  # show recently played tracks
sp queue                   # show the current queue
sp queue add spotify:track:4iV5W9uYEdYUVa79Axb7Rh
sp volume 80               # set volume to 80%
```

### Search

```sh
sp search tracks "dark side of the moon"
sp search albums --mine    # browse your saved albums
sp search playlists --mine
```

### Shelf

```sh
sp shelf add lofi spotify:playlist:37i9dQZF1DX3Ogo9pFvBkY
sp shelf                   # list all shortcuts
sp play lofi               # play a shelf item instantly
sp shelf drop lofi         # remove a shortcut
```

### Devices

```sh
sp devices                 # list available devices
sp activate "MacBook Pro"  # switch playback device
sp open                    # open Spotify on this device
sp close                   # close Spotify on this device
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
