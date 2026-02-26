BINARY     := fediscord
CMD_PATH   := ./cmd/fediscord
BUILD_DIR  := build
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -ldflags "-X main.version=$(VERSION) -s -w"

.PHONY: all build clean install uninstall deps tidy vet

all: build

deps:
	go mod tidy
	go get golang.org/x/term

tidy:
	go mod tidy

vet:
	go vet ./...

build: deps vet
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) $(CMD_PATH)

install: build
	install -Dm755 $(BUILD_DIR)/$(BINARY) $(DESTDIR)/usr/local/bin/$(BINARY)

uninstall:
	rm -f $(DESTDIR)/usr/local/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)
