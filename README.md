# fediscord

A terminal-based utility for establishing a verified connection between a Mastodon API-compatible Fediverse account and a Discord user profile. The tool interacts with Discord's internal connections API to generate the authorisation URL required to confirm Fediverse account ownership from within Discord. All credential management is performed locally, with support for GPG symmetric encryption on Linux and macOS.

---

## Table of Contents

- [Overview](#overview)
- [Platform Compatibility](#platform-compatibility)
- [Supported Operating Systems and Architectures](#supported-operating-systems-and-architectures)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Building](#building)
  - [Current Host Platform](#current-host-platform)
  - [Cross-Compilation](#cross-compilation)
  - [Release Build (All Platforms)](#release-build-all-platforms)
- [Installation](#installation)
- [Uninstallation](#uninstallation)
- [Usage](#usage)
  - [Main Menu](#main-menu)
  - [1 — Set Up Configuration](#1--set-up-configuration)
  - [2 — Generate Connection URL](#2--generate-connection-url)
  - [3 — View Stored Configuration](#3--view-stored-configuration)
  - [4 — Update Discord Token](#4--update-discord-token)
  - [5 — Update Fediverse Handle](#5--update-fediverse-handle)
  - [6 — Change Encryption Settings](#6--change-encryption-settings)
  - [7 — Delete All Data](#7--delete-all-data)
  - [8 — Exit](#8--exit)
- [Token Security](#token-security)
- [Encryption](#encryption)
- [Configuration Storage Paths](#configuration-storage-paths)
- [Makefile Reference](#makefile-reference)
- [Windows Considerations](#windows-considerations)
- [macOS Considerations](#macos-considerations)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Licence](#licence)

---

## Overview

`fediscord` is a portable command-line tool implemented in Go. Its primary function is to facilitate the linkage of a Fediverse identity to a Discord account by generating the authorisation URL that the Discord connections system requires. The tool communicates with the Discord API v9 endpoint designated for Mastodon-type connections (`/api/v9/connections/mastodon/authorize`).

The entire operation is conducted through an interactive terminal interface. No positional arguments or flag-based invocation are required; all inputs are solicited at runtime through a numbered menu system. Credentials are persisted in a platform-appropriate local directory with restrictive access permissions, and the Discord token may optionally be protected at rest using GPG symmetric AES-256 encryption on supported platforms.

The application has been designed to operate across Linux distributions, macOS (both Intel and Apple Silicon), and Microsoft Windows, with platform-specific behaviour abstracted through Go build constraints and a dedicated terminal abstraction package.

---

## Platform Compatibility

The tool is designed for use with Fediverse instances that implement the Mastodon-compatible REST API. Compatibility is determined by the presence and structure of the `/api/v1/instance` response returned by the instance server. The following platforms are known to be compatible:

| Platform  | API Compatibility |
|-----------|-------------------|
| Mastodon  | Supported         |
| Akkoma    | Supported         |
| Pleroma   | Supported         |
| GlitchSoc | Supported         |
| Hometown  | Supported         |

The following platforms implement a divergent API architecture (the Misskey API) and are consequently incompatible with the Discord Mastodon connection mechanism:

| Platform  | API Compatibility |
|-----------|-------------------|
| Misskey   | Incompatible      |
| Firefish  | Incompatible      |
| Calckey   | Incompatible      |
| Foundkey  | Incompatible      |

Instance compatibility is automatically verified during the set-up procedure by querying the instance's API endpoint directly. The tool will report the result and, in the event of an incompatible or unreachable instance, will prompt the user to confirm whether they wish to proceed.

---

## Supported Operating Systems and Architectures

The following target combinations are supported and may be compiled via the Makefile:

| Operating System | Architecture     | Binary Suffix              |
|------------------|------------------|----------------------------|
| Linux            | x86\_64 (amd64)  | `fediscord-linux-amd64`    |
| Linux            | ARM64 (AArch64)  | `fediscord-linux-arm64`    |
| Linux            | ARMv6            | `fediscord-linux-arm`      |
| Linux            | x86 (32-bit)     | `fediscord-linux-386`      |
| macOS            | Intel (amd64)    | `fediscord-darwin-amd64`   |
| macOS            | Apple Silicon    | `fediscord-darwin-arm64`   |
| Windows          | x86\_64 (amd64)  | `fediscord-windows-amd64.exe` |
| Windows          | ARM64            | `fediscord-windows-arm64.exe` |
| Windows          | x86 (32-bit)     | `fediscord-windows-386.exe`   |

---

## Project Structure

```
fediscord/
├── cmd/
│   └── fediscord/
│       ├── main.go       Entry point, menu loop, and application lifecycle management
│       ├── setup.go      Configuration set-up and configuration view handlers
│       └── actions.go    URL generation, credential update, encryption, and deletion handlers
├── pkg/
│   ├── config/
│   │   └── config.go     Platform-aware configuration path resolution and directory initialisation
│   ├── discord/
│   │   └── discord.go    Discord API v9 interaction and authorisation URL construction
│   ├── fediverse/
│   │   └── fediverse.go  Fediverse handle validation and Mastodon API instance verification
│   ├── storage/
│   │   └── storage.go    Credential persistence with optional GPG encryption support
│   ├── terminal/
│   │   ├── terminal_unix.go     Unix/Linux/macOS terminal operations (build-constrained)
│   │   └── terminal_windows.go  Windows terminal operations (build-constrained)
│   └── ui/
│       └── ui.go         Terminal input/output helpers, menu rendering, and user prompts
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

The `cmd` layer is maintained as a thin orchestration layer that delegates all substantive logic to the packages under `pkg`. This separation ensures that individual packages may be imported independently or tested in isolation. Platform-specific behaviour is encapsulated within the `terminal` package, which employs Go build constraints (`//go:build`) to select the appropriate implementation at compile time.

---

## Requirements

The following dependencies are required to build and operate this tool:

- **Go 1.21 or later** — The Go toolchain is required to compile the project. The official distribution is available at [https://go.dev/dl](https://go.dev/dl).
- **make** — Required to utilise the provided `Makefile`. This utility is installed by default on Linux and macOS. On Windows, it may be obtained via [GnuWin32](http://gnuwin32.sourceforge.net/packages/make.htm), [Chocolatey](https://chocolatey.org/packages/make), or through the Windows Subsystem for Linux (WSL).
- **gpg** *(optional; Linux and macOS only)* — Required to enable GPG-based encrypted token storage. GPG is not supported on Windows; tokens on that platform are stored in plain text with restricted file-system permissions. Installation instructions by platform:
  - Ubuntu / Debian: `sudo apt install gnupg`
  - Fedora / RHEL: `sudo dnf install gnupg2`
  - Arch Linux: `sudo pacman -S gnupg`
  - macOS (Homebrew): `brew install gnupg`
- **git** *(optional)* — Utilised by the Makefile to derive version information from git tags via `git describe`. If git is unavailable, the version string defaults to `dev`.

---

## Building

### Current Host Platform

To compile the binary for the operating system and architecture of the current build host:

```sh
make build
```

This target executes the following steps in sequence:

1. Fetches and tidies all Go module dependencies (`golang.org/x/term`).
2. Executes `go vet` across all packages to identify potential static analysis issues.
3. Compiles the binary and places it at `./fediscord` (or `./fediscord.exe` on Windows when using `go build` directly).

The version string embedded in the binary is derived from `git describe --tags --always --dirty`. In the absence of git tags, the string `dev` is substituted.

To verify the build was successful:

```sh
./fediscord
```

The interactive menu will be presented immediately upon execution.

### Cross-Compilation

Individual platform targets may be compiled without executing a full release build. Each target is self-contained and places the resulting binary under `./dist/`:

```sh
make linux-amd64
make linux-arm64
make linux-arm
make linux-386
make darwin-amd64
make darwin-arm64
make windows-amd64
make windows-arm64
make windows-386
```

No external cross-compilation toolchain is required. Go's native cross-compilation support is used throughout; the `GOOS` and `GOARCH` environment variables are set explicitly for each target.

### Release Build (All Platforms)

To compile binaries for all supported platforms simultaneously:

```sh
make release
```

Upon completion, all binaries will be present under `./dist/` with platform-indicative suffixes. A directory listing is printed automatically at the conclusion of the build.

---

## Installation

To install the binary to the system-wide executable path:

```sh
make install
```

This installs `fediscord` to `/usr/local/bin/fediscord` using `install -Dm755`. Elevated privileges (`sudo`) may be required depending on the write permissions of the target directory.

To install to a custom path prefix:

```sh
make install DESTDIR=/home/youruser/.local
```

This places the binary at `/home/youruser/.local/usr/local/bin/fediscord`. Ensure this path is present in your `$PATH` environment variable.

**On Windows**, it is recommended to place the compiled binary in a directory that is registered in the system or user `PATH` environment variable, such as `C:\Users\YourName\bin\`. The `make install` target is not applicable on Windows without a Unix-compatible environment such as WSL or MSYS2.

---

## Uninstallation

```sh
make uninstall
```

This removes the binary at `/usr/local/bin/fediscord`. Supply `DESTDIR` if a custom prefix was used during installation:

```sh
make uninstall DESTDIR=/home/youruser/.local
```

This operation does **not** remove locally stored configuration data. To remove all stored credentials and settings, use option `7` from within the tool, or delete the configuration directory manually. Refer to [Configuration Storage Paths](#configuration-storage-paths) for the applicable path on your operating system.

---

## Usage

Execute the binary from any location in the terminal:

```sh
fediscord
```

Or from the project root prior to installation:

```sh
./fediscord
```

No command-line arguments, flags, or environment variables are required. All interactions are conducted through the numbered interactive menu.

### Main Menu

```
╔═══════════════════════════════════════════════════════════╗
║  Fediverse to Discord Connection Tool (Mastodon API)     ║
╠═══════════════════════════════════════════════════════════╣
║  Main Menu                                               ║
╚═══════════════════════════════════════════════════════════╝

1) Set Up Configuration (Discord Token + Fediverse Handle)
2) Generate Connection URL
3) View Stored Configuration
4) Update Discord Token
5) Update Fediverse Handle
6) Change Encryption Settings
7) Delete All Data
8) Exit

───────────────────────────────────────────────────────────
Compatible Platforms (Mastodon API):
  + Mastodon  + Akkoma  + Pleroma  + GlitchSoc  + Hometown
Incompatible Platforms:
  - Misskey  - Firefish  - Calckey  - Foundkey
───────────────────────────────────────────────────────────
```

---

### 1 — Set Up Configuration

This is the initial configuration procedure and must be completed before other functions of the tool can be utilised. The procedure consists of two sequential steps.

**Step 1: Discord Token**

The user is prompted to provide their Discord user-level account token. This token is distinct from a Discord bot token and is the credential used by the Discord client application itself. A reference guide for retrieving this token is displayed within the tool. Prior to token input, the user is asked to select a storage method:

- **Encrypted (GPG AES-256):** Available on Linux and macOS where GPG is installed. The token is encrypted symmetrically using a passphrase supplied by the user. The passphrase will be required on subsequent accesses.
- **Plain text:** The token is stored without encryption, with `0600` file-system permissions. This option is the only available method on Windows.

**Step 2: Fediverse Handle**

The user is prompted to provide their Fediverse handle in the format `username@instance.domain` (the leading `@` is optional and will be stripped automatically). The tool performs the following validation steps:

1. The handle is validated against a regular expression to confirm structural conformance.
2. The instance domain is extracted from the handle.
3. A live HTTP request is made to `https://instance.domain/api/v1/instance` to verify that the instance is reachable and operating on a Mastodon API-compatible platform.
4. The version string returned by the instance is examined for known incompatible platform identifiers.

In the event that the API check fails, the user is given the option to proceed regardless.

Upon successful completion of both steps, the token and handle are persisted to the configuration directory with `0600` permissions.

---

### 2 — Generate Connection URL

This function retrieves the stored token and handle, then issues a GET request to the Discord API endpoint:

```
https://discord.com/api/v9/connections/mastodon/authorize?handle=@username@instance.domain
```

The `authorization` header is populated with the stored Discord token. If the request is successful, Discord returns a JSON object containing an authorisation URL, which is displayed in the terminal.

**To complete the connection process:**

1. Copy the authorisation URL displayed in the terminal.
2. Paste it into a web browser.
3. If not already authenticated, log in to your Fediverse account on the relevant instance.
4. Authorise the connection request presented by the OAuth flow.
5. Discord will record the verified Fediverse account link upon completion.

If the request fails, the tool presents an error message along with a list of commonly applicable diagnostic considerations.

---

### 3 — View Stored Configuration

Displays a summary of the currently stored configuration, comprising:

- The presence or absence of a stored Discord token.
- The storage method in use (GPG-encrypted or plain text).
- A 10-character prefix preview of the token (the full value is not displayed).
- The stored Fediverse handle and its extracted instance domain.

No sensitive data beyond the preview prefix is rendered to the terminal.

---

### 4 — Update Discord Token

Replaces the currently stored Discord token with a new value supplied by the user. The existing encryption preference is retained; the replacement token will be stored using the same method as its predecessor.

---

### 5 — Update Fediverse Handle

Replaces the currently stored Fediverse handle. The replacement handle undergoes the same validation and instance API check procedure as described in Section 1, Step 2.

---

### 6 — Change Encryption Settings

Permits modification of the token storage method. If a token is already stored, it is retrieved and re-stored using the newly selected method. The transition is bidirectional:

- **Plain text to encrypted:** The plain-text token is read, encrypted via GPG, written to the encrypted token file, and the plain-text file is removed.
- **Encrypted to plain text:** The encrypted token is decrypted via GPG, written to the plain-text file, and the encrypted file is removed.

On Windows, this option will default to plain-text storage as GPG is not available on that platform.

---

### 7 — Delete All Data

Permanently removes the entire configuration directory and all contents, including the stored token, handle, and encryption preference. The user must type `DELETE` (all uppercase, case-sensitive) to confirm the operation.

**This action is irreversible.** The configuration directory must be recreated through the set-up procedure (Option 1) if the tool is to be used again following deletion.

---

### 8 — Exit

Clears the terminal and terminates the process.

---

## Token Security

A Discord user token is a credential of the highest sensitivity. It is functionally equivalent to a username and password combination in that its possession grants unrestricted access to the associated Discord account, including the ability to read private messages, modify account settings, and perform any action that the account owner is authorised to perform.

The following practices are strongly recommended:

- Under no circumstances share the token with any person, application, or service beyond this tool.
- Do not paste the token in any publicly accessible location, including chat platforms, public code repositories, or issue trackers.
- Do not store the token in locations accessible to other users of the same system unless encryption is employed.
- If the token is believed to have been exposed, change the Discord account password immediately. This action invalidates the current token and necessitates the generation of a new one.

`fediscord` stores the token at the path determined by [Configuration Storage Paths](#configuration-storage-paths) with `0600` file-system permissions, ensuring that the file is readable only by the owning user account.

---

## Encryption

When GPG is available on the system (Linux and macOS only) and the user selects the encrypted storage method, `fediscord` invokes the following command to encrypt the token:

```sh
gpg --symmetric --cipher-algo AES256 --output <token_path>
```

The token is passed to `gpg` via standard input to avoid exposure in process argument lists. GPG will prompt for a passphrase, which is required each time the token must be decrypted (e.g. when generating a connection URL or updating the token). The cipher employed is AES-256, a symmetric block cipher widely accepted as suitable for the protection of sensitive data at rest.

On Windows, GPG integration is not available. Tokens on this platform are stored in plain text within the configuration directory, relying on the operating system's file-system access controls for protection. Users who require enhanced security on Windows should consider operating within WSL (Windows Subsystem for Linux), where GPG is available.

---

## Configuration Storage Paths

The configuration directory is determined at runtime based on the operating system:

| Operating System | Configuration Directory                                       |
|------------------|---------------------------------------------------------------|
| Linux            | `~/.fediverse-discord/`                                       |
| macOS            | `~/Library/Application Support/fediverse-discord/`           |
| Windows          | `%APPDATA%\fediverse-discord\`                                |

The following files are maintained within this directory:

| File                    | Purpose                                                 |
|-------------------------|---------------------------------------------------------|
| `discord_token.txt`     | Discord token stored in plain text                      |
| `discord_token.enc`     | Discord token stored GPG-encrypted (AES-256)            |
| `fediverse_handle.txt`  | Stored Fediverse handle (`username@instance.domain`)    |
| `.use_encryption`       | Encryption preference flag (`true` or `false`)          |

The configuration directory is created with `0700` permissions; individual files are created with `0600` permissions. Only one of `discord_token.txt` or `discord_token.enc` will be present at any given time; a change in storage method results in the removal of the superseded file.

---

## Makefile Reference

Executing `make` without specifying a target displays the full help output, including all available targets and their descriptions.

| Target          | Description                                                         |
|-----------------|---------------------------------------------------------------------|
| `build`         | Compile the binary for the host platform; output: `./fediscord`     |
| `install`       | Install the binary to `$(DESTDIR)/usr/local/bin/fediscord`          |
| `uninstall`     | Remove the binary from `$(DESTDIR)/usr/local/bin/fediscord`         |
| `clean`         | Remove `./fediscord` and the `./dist/` directory                    |
| `deps`          | Execute `go mod tidy` and fetch `golang.org/x/term`                 |
| `tidy`          | Execute `go mod tidy`                                               |
| `vet`           | Execute `go vet ./...` across all packages                          |
| `release`       | Compile binaries for all supported platforms to `./dist/`           |
| `linux-amd64`   | Compile for Linux x86\_64                                           |
| `linux-arm64`   | Compile for Linux ARM64                                             |
| `linux-arm`     | Compile for Linux ARMv6                                             |
| `linux-386`     | Compile for Linux x86 (32-bit)                                      |
| `darwin-amd64`  | Compile for macOS Intel                                             |
| `darwin-arm64`  | Compile for macOS Apple Silicon                                     |
| `windows-amd64` | Compile for Windows x86\_64; output: `.exe`                         |
| `windows-arm64` | Compile for Windows ARM64; output: `.exe`                           |
| `windows-386`   | Compile for Windows x86 (32-bit); output: `.exe`                    |

---

## Windows Considerations

The following platform-specific notes apply to operation on Microsoft Windows:

- **GPG encryption is not supported.** The Discord token will be stored in plain text within the `%APPDATA%\fediverse-discord\` directory. File-system access controls native to Windows (NTFS permissions) provide the primary means of access restriction.
- **Terminal clear screen behaviour differs.** On Windows, the ANSI escape sequence used to clear the terminal on Unix-like systems is not emitted. The screen is not cleared between menu transitions on Windows terminals that do not support ANSI. This does not affect functionality.
- **The `make install` target is not supported natively on Windows.** To install the binary system-wide, manually copy the compiled `.exe` file to a directory present in your `PATH` environment variable.
- **Unicode box-drawing characters** used in the menu interface require a terminal emulator with appropriate Unicode support, such as Windows Terminal. The tool functions correctly in Windows Terminal; compatibility with the legacy `cmd.exe` console host is not guaranteed for all visual elements.
- **Compilation on Windows** requires either Go for Windows, WSL, or MSYS2 with make installed.

---

## macOS Considerations

- The configuration directory is located at `~/Library/Application Support/fediverse-discord/`, consistent with macOS application data conventions.
- GPG encryption is fully supported. Install GPG via Homebrew (`brew install gnupg`) to enable encrypted token storage.
- Both Intel (`darwin/amd64`) and Apple Silicon (`darwin/arm64`) targets are supported. The `darwin-arm64` target produces a native binary for M-series hardware and should be preferred on those systems.
- The binary is unsigned. On macOS 12 (Monterey) and later, Gatekeeper may prevent execution of unsigned binaries downloaded from the internet. To permit execution after download, run: `xattr -d com.apple.quarantine ./fediscord`.

---

## Troubleshooting

**`The requested credential was not found in the local configuration store`**

No token or handle has been configured. Complete the initial set-up procedure using Option 1 from the main menu.

**`GPG is not available on this system, however an encrypted token was detected`**

GPG was removed or is no longer accessible in `$PATH` after an encrypted token was stored. Reinstall GPG and retry, or use Option 6 to convert the token to plain-text storage prior to removing GPG.

**`The instance did not return a valid Mastodon API v1 response`**

The specified instance is either unreachable over the network, returning an unexpected response format, or is operating on a platform that does not implement the Mastodon v1 API (e.g. Misskey). Refer to [Platform Compatibility](#platform-compatibility) and verify that the instance domain is correct.

**`The Discord API did not return a valid authorisation URL`**

The Discord token is most likely invalid or has been revoked. Retrieve a new token and update it using Option 4. Additionally verify that network connectivity to `discord.com` is available.

**`The token encryption process failed`**

GPG encountered an error during the symmetric encryption step. Common causes include: the user cancelled the passphrase prompt, an empty passphrase was entered, or the GPG agent is in an inconsistent state. To diagnose, execute `gpg --symmetric --cipher-algo AES256 /tmp/test_file` in the terminal directly.

**Gatekeeper blocks execution on macOS**

Remove the quarantine extended attribute: `xattr -d com.apple.quarantine ./fediscord`

**Build fails with `go: command not found`**

The Go toolchain is not installed or is not present in `$PATH`. Refer to [https://go.dev/dl](https://go.dev/dl) for installation instructions appropriate to your platform.

---

## Contributing

Contributions are welcomed in the form of bug reports, feature requests, and pull requests. Please submit these through the GitHub repository:

[https://github.com/jimed-rand/fediscord](https://github.com/jimed-rand/fediscord)

When contributing code, the following conventions must be observed:

- Inline comments must not be added to source files.
- Hardcoded values are not permitted; all configurable values must be externalised through variables or the configuration system.
- Platform-specific behaviour must be implemented using Go build constraints within the `pkg/terminal` package or an analogous dedicated package, rather than through runtime `if runtime.GOOS` checks dispersed throughout the codebase.
- The structural separation between `cmd/` (orchestration) and `pkg/` (reusable logic) must be maintained.
- All submissions must pass `go vet ./...` without findings.

---

## Licence

GNU General Public License v2.0 — James "Jim" Ed Randson

See [LICENSE](./LICENSE) for the full licence text.
