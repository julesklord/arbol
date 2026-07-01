BINARY=arbol
VERSION ?= $(shell cat VERSION 2>/dev/null || echo dev)
LDFLAGS=-s -w -X main.version=$(VERSION)

.PHONY: build build-all install clean test release-snapshot

# Build: prefer Go build if cmd/arbol/main.go exists, otherwise no-op
build:
	@if [ -f "cmd/arbol/main.go" ]; then \
		go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/arbol; \
		echo "Built $(BINARY) v$(VERSION) (Go version)"; \
	else \
		echo "No Go implementation found — nothing to build"; \
	fi

# Cross-compile for all supported platforms
BUILD_DIR=dist
BUILD_DIRS=$(BUILD_DIR)/linux-amd64 $(BUILD_DIR)/linux-arm64 $(BUILD_DIR)/darwin-amd64 $(BUILD_DIR)/darwin-arm64

build-all: $(BUILD_DIRS)
	@echo "Building arbol v$(VERSION) for all platforms..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/linux-amd64/$(BINARY) ./cmd/arbol
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/linux-arm64/$(BINARY) ./cmd/arbol
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/darwin-amd64/$(BINARY) ./cmd/arbol
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/darwin-arm64/$(BINARY) ./cmd/arbol
	@echo "Built all platforms"

$(BUILD_DIRS):
	mkdir -p $@

# Install: install built binary and assets
install: build
	@if [ -d "ascii" ]; then \
		mkdir -p /usr/local/share/arbol/ascii; \
		cp -r ascii/* /usr/local/share/arbol/ascii/; \
		echo "Installed ASCII assets to /usr/local/share/arbol/ascii/"; \
	fi
	@if [ -f "$(BINARY)" ]; then \
		install -m 0755 $(BINARY) /usr/local/bin/$(BINARY); \
		echo "Installed built $(BINARY) binary to /usr/local/bin/$(BINARY)"; \
	else \
		echo "Nothing to install (binary not found)"; exit 1; \
	fi


# Test: runs the Go test suite
test: build
	go test -v ./cmd/arbol

release-snapshot:
	@goreleaser release --snapshot --clean

clean:
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)


