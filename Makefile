VERSION := $(shell git describe --tags --dirty --always)
UNAME := $(shell uname)

# Config directory (database): uses os.UserConfigDir() pattern
# - macOS: ~/Library/Application Support/Footprint
# - Linux: $XDG_CONFIG_HOME/Footprint or ~/.config/Footprint
ifeq ($(UNAME),Darwin)
	FP_CONFIG_DIR := $(HOME)/Library/Application Support/Footprint
else
	FP_CONFIG_DIR := $(or $(XDG_CONFIG_HOME),$(HOME)/.config)/Footprint
endif

# Data directory (exports): uses XDG_DATA_HOME pattern
# - macOS: ~/Library/Application Support/footprint
# - Linux: $XDG_DATA_HOME/footprint or ~/.local/share/footprint
ifeq ($(UNAME),Darwin)
	FP_DATA_DIR := $(HOME)/Library/Application Support/footprint
else
	FP_DATA_DIR := $(or $(XDG_DATA_HOME),$(HOME)/.local/share)/footprint
endif

.PHONY: build test test-nocache test-actions test-hooks test-export test-backfill wipe

build: test
	go build \
		-ldflags "-X github.com/Skryensya/footprint/internal/app.Version=$(VERSION)" \
		-o fp \
		./cmd/fp

test:
	go test ./...

test-nocache:
	go test -count=1 ./...

test-actions:
	go test ./internal/actions

test-hooks: build
	./scripts/test-hooks.sh

test-export: build
	./scripts/test-export-flow.sh

test-backfill: build
	./scripts/test-backfill.sh

wipe:
	rm -rf "$(FP_CONFIG_DIR)"
	rm -rf "$(FP_DATA_DIR)"
	rm -f ~/.fprc
	@echo "Wiped database, exports, and config"
