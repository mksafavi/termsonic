# termsonic - a TUI Subsonic client

This project implements a terminal-based client for any [Subsonic](https://www.subsonic.org)-compatible server.

## Building

This application requires [Go](https://go.dev) version 1.19 at minimum.

```
$ git clone https://git.sixfoisneuf.fr/termsonic && cd termsonic
$ go build -o termsonic ./cmd
```

## Configuration

The application reads its configuration from `$XDG_CONFIG_DIR/termsonic.toml`, or `~/.config/termsonic.toml` if `XDG_CONFIG_DIR` doesn't exist.

On Windows, it reads its configuration from `%APPDATA%\\Termsonic\\termsonic.toml`.

You can edit the configuration from inside the app, or by passing parameters on the command line (see `--help`).