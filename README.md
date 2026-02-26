# fediscord

A terminal tool to connect Mastodon API-compatible Fediverse accounts to Discord.

**Compatible platforms:** Mastodon, Akkoma, Pleroma, GlitchSoc, Hometown  
**Incompatible platforms:** Misskey, Firefish, Calckey, Foundkey

## Requirements

- Go 1.21 or later
- `gpg` (optional, for encrypted token storage)

## Build

```sh
make build
```

The binary will be placed in `build/fediscord`.

## Install

```sh
make install
```

Installs to `/usr/local/bin/fediscord`. Override with `DESTDIR`:

```sh
make install DESTDIR=/custom/prefix
```

## Uninstall

```sh
make uninstall
```

## Usage

Run the binary and follow the interactive menu:

```
1) Setup Configuration (Discord Token + Fediverse Handle)
2) Generate Connection URL
3) View Stored Configuration
4) Update Discord Token
5) Update Fediverse Handle
6) Change Encryption Settings
7) Delete All Data
8) Exit
```

Configuration is stored in `~/.fediverse-discord/`. The Discord token can be stored either GPG-encrypted (AES256) or in plain text, depending on user choice and GPG availability.

## Project Structure

```
fediscord/
├── cmd/
│   └── fediscord/
│       ├── main.go        # Entry point and main menu loop
│       └── actions.go     # Menu action handlers
├── internal/
│   ├── config/
│   │   └── config.go      # Configuration paths and directory init
│   ├── discord/
│   │   └── discord.go     # Discord API interaction
│   ├── fediverse/
│   │   └── fediverse.go   # Handle validation and instance checking
│   ├── storage/
│   │   └── storage.go     # Token and handle read/write with GPG support
│   └── ui/
│       └── ui.go          # Terminal input/output helpers
├── go.mod
├── Makefile
└── README.md
```

## License

MIT — James "Jim" Ed Randson
