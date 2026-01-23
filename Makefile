VERSION := $(shell git describe --tags --dirty --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/footprint-tools/footprint-cli/internal/app.Version=$(VERSION)"

.PHONY: all build test lint fmt clean install wipe integration

# Default target
all: build

# Build binary (runs tests first)
build: test
	go build $(LDFLAGS) -o fp ./cmd/fp

# Build without tests (for quick iteration)
build-fast:
	go build $(LDFLAGS) -o fp ./cmd/fp

# Run unit tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	rm -f fp
	go clean

# Install to GOPATH/bin
install: test
	go install $(LDFLAGS) ./cmd/fp

# Run integration tests (slow, requires built binary)
integration: build
	./scripts/test-hooks.sh
	./scripts/test-export-flow.sh
	./scripts/test-backfill.sh

# Wipe all local data (database, exports, config)
wipe:
	rm -rf "$(HOME)/Library/Application Support/Footprint"
	rm -rf "$(HOME)/Library/Application Support/footprint"
	rm -rf "$${XDG_CONFIG_HOME:-$$HOME/.config}/Footprint"
	rm -rf "$${XDG_DATA_HOME:-$$HOME/.local/share}/footprint"
	rm -f ~/.fprc
	@echo "Wiped database, exports, and config"
