BINARY     := fediscord
CMD_PATH   := ./cmd/fediscord
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -ldflags "-X main.version=$(VERSION) -s -w"

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	linux/arm \
	linux/386 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64 \
	windows/386

.PHONY: all build install uninstall clean deps tidy vet release help \
	linux-amd64 linux-arm64 linux-arm linux-386 \
	darwin-amd64 darwin-arm64 \
	windows-amd64 windows-arm64 windows-386

all: help

help:
	@printf "\n"
	@printf "  fediscord — Fediverse to Discord Connection Tool\n"
	@printf "  Version : %s\n\n" "$(VERSION)"
	@printf "  ── Standard Targets ────────────────────────────────────────\n\n"
	@printf "    make build          Compile for the current host platform\n"
	@printf "    make install        Install to $(DESTDIR)/usr/local/bin\n"
	@printf "    make uninstall      Remove from $(DESTDIR)/usr/local/bin\n"
	@printf "    make clean          Remove all compiled artefacts\n"
	@printf "    make deps           Fetch and tidy module dependencies\n"
	@printf "    make tidy           Run go mod tidy\n"
	@printf "    make vet            Run static analysis via go vet\n"
	@printf "    make release        Compile all supported platforms\n\n"
	@printf "  ── Platform-Specific Targets ───────────────────────────────\n\n"
	@printf "    make linux-amd64    Linux   — x86_64\n"
	@printf "    make linux-arm64    Linux   — ARM64 (AArch64)\n"
	@printf "    make linux-arm      Linux   — ARMv6 (Raspberry Pi etc.)\n"
	@printf "    make linux-386      Linux   — x86 (32-bit)\n"
	@printf "    make darwin-amd64   macOS   — Intel\n"
	@printf "    make darwin-arm64   macOS   — Apple Silicon (M-series)\n"
	@printf "    make windows-amd64  Windows — x86_64\n"
	@printf "    make windows-arm64  Windows — ARM64\n"
	@printf "    make windows-386    Windows — x86 (32-bit)\n\n"
	@printf "  ── Override ────────────────────────────────────────────────\n\n"
	@printf "    make install DESTDIR=/custom/prefix\n\n"

deps:
	@printf "  [--] Fetching module dependencies...\n"
	go mod tidy
	go get golang.org/x/term
	@printf "  [OK] Dependencies are up to date.\n"

tidy:
	@printf "  [--] Executing go mod tidy...\n"
	go mod tidy
	@printf "  [OK] Module graph has been tidied.\n"

vet:
	@printf "  [--] Executing static analysis (go vet)...\n"
	go vet ./...
	@printf "  [OK] Static analysis completed without findings.\n"

build: deps vet
	@printf "  [--] Compiling %s for host platform (version: %s)...\n" "$(BINARY)" "$(VERSION)"
	go build $(LDFLAGS) -o ./$(BINARY) $(CMD_PATH)
	@printf "  [OK] Binary produced: ./%s\n" "$(BINARY)"

install: build
	@printf "  [--] Installing %s to %s/usr/local/bin...\n" "$(BINARY)" "$(DESTDIR)"
	install -Dm755 ./$(BINARY) $(DESTDIR)/usr/local/bin/$(BINARY)
	@printf "  [OK] Installation complete.\n"

uninstall:
	@printf "  [--] Removing %s from %s/usr/local/bin...\n" "$(BINARY)" "$(DESTDIR)"
	rm -f $(DESTDIR)/usr/local/bin/$(BINARY)
	@printf "  [OK] Uninstall complete.\n"

clean:
	@printf "  [--] Removing compiled artefacts...\n"
	rm -f ./$(BINARY)
	rm -rf ./dist
	@printf "  [OK] Working directory is clean.\n"

dist:
	mkdir -p dist

linux-amd64: dist deps vet
	@printf "  [--] Compiling for Linux / amd64...\n"
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-linux-amd64\n"

linux-arm64: dist deps vet
	@printf "  [--] Compiling for Linux / arm64...\n"
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-linux-arm64\n"

linux-arm: dist deps vet
	@printf "  [--] Compiling for Linux / arm (ARMv6)...\n"
	GOOS=linux GOARCH=arm GOARM=6 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-linux-arm\n"

linux-386: dist deps vet
	@printf "  [--] Compiling for Linux / 386 (32-bit)...\n"
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o dist/$(BINARY)-linux-386 $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-linux-386\n"

darwin-amd64: dist deps vet
	@printf "  [--] Compiling for macOS / amd64 (Intel)...\n"
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-darwin-amd64\n"

darwin-arm64: dist deps vet
	@printf "  [--] Compiling for macOS / arm64 (Apple Silicon)...\n"
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-darwin-arm64\n"

windows-amd64: dist deps vet
	@printf "  [--] Compiling for Windows / amd64...\n"
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-windows-amd64.exe\n"

windows-arm64: dist deps vet
	@printf "  [--] Compiling for Windows / arm64...\n"
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-arm64.exe $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-windows-arm64.exe\n"

windows-386: dist deps vet
	@printf "  [--] Compiling for Windows / 386 (32-bit)...\n"
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o dist/$(BINARY)-windows-386.exe $(CMD_PATH)
	@printf "  [OK] dist/$(BINARY)-windows-386.exe\n"

release: linux-amd64 linux-arm64 linux-arm linux-386 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64 windows-386
	@printf "\n  [OK] All platform binaries have been produced under ./dist/\n\n"
	@ls -lh dist/
