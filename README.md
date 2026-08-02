<div align="center">

<img width="400" src="etc/img.svg" alt="sp demo" />

**A fast, minimal Spotify CLI and MCP server**

[![release](https://img.shields.io/github/v/release/matthew-collett/sp?color=059669)](https://github.com/matthew-collett/sp/releases)
[![ci](https://img.shields.io/github/actions/workflow/status/matthew-collett/sp/ci.yaml)](https://github.com/matthew-collett/sp/actions)
[![mcp](https://img.shields.io/badge/mcp-compatible-F87171)](https://modelcontextprotocol.io)
[![license](https://img.shields.io/github/license/matthew-collett/sp?color=F97316)](LICENSE)

Control Spotify from your terminal or your LLM: search, queue, devices, playback, and quick-play shortcuts, from wherever you work.

<img src="etc/demo.gif" alt="sp demo" />

</div>

## Installation

### macOS / Linux (Homebrew)

```sh
brew tap matthew-collett/sp
brew trust matthew-collett/sp
brew install sp
```

### macOS / Linux (curl)

```sh
curl -fsSL https://raw.githubusercontent.com/matthew-collett/sp/main/tools/install.sh | sh
```

Installs both `sp` and `sp-mcp` to `/usr/local/bin`.

For other platforms, download pre-built binaries from the [releases page](https://github.com/matthew-collett/sp/releases).

## Setup

`sp` uses the Spotify Web API and requires a developer application.

Playback control (play, pause, skip, queue, volume, etc.) requires a Spotify Premium account. Free-tier accounts will get 403 errors from those endpoints.

1. Create an app at the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Add `http://localhost:8080/callback` as a redirect URI
3. Copy your Client ID and Client Secret

```sh
sp configure   # enter your credentials
sp login       # authenticate via browser
```

## MCP Server

`sp-mcp` brings the same controls to any [MCP](https://modelcontextprotocol.io)-compatible LLM as native tools. Ask naturally instead of typing commands:

- "What's playing?"
- "Put on my lofi from my shelf"
- "Skip this and like the next one"
- "Turn shuffle on and set volume to 60"
- "What have I been listening to recently?"

### Claude Code

```sh
claude mcp add sp --scope user -- sp-mcp
```

### Other clients

For any other MCP-compatible client, add this to your config file:

```json
{
  "mcpServers": {
    "sp": {
      "command": "sp-mcp"
    }
  }
}
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
sp queue add spotify:...   # add to current queue
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
sp shelf add lofi spotify:... # add a shortcut
sp shelf                      # list all shortcuts
sp play lofi                  # play a shelf item instantly
sp shelf drop lofi            # remove a shortcut
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
