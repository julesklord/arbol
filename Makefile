BINARY=mini-fetch

.PHONY: build install clean

# Build: prefer Go build if main.go exists, otherwise no-op
build:
	@if [ -f "main.go" ]; then \
		go build -o $(BINARY) .; \
		echo "Built $(BINARY)"; \
	else \
		echo "No main.go found — nothing to build"; \
	fi

# Install: prefer the built binary if present, otherwise install script
install: build
	@if [ -f "$(BINARY)" ]; then \
		install -m 0755 $(BINARY) /usr/local/bin/$(BINARY); \
		echo "Installed /usr/local/bin/$(BINARY)"; \
	elif [ -f "scripts/mini-fetch.sh" ]; then \
		install -m 0755 scripts/mini-fetch.sh /usr/local/bin/$(BINARY); \
		echo "Installed scripts/mini-fetch.sh as /usr/local/bin/$(BINARY)"; \
	else \
		echo "Nothing to install"; exit 1; \
	fi

clean:
	rm -f $(BINARY)
