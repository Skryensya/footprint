VERSION := $(shell git describe --tags --dirty --always)

.PHONY: build test test-actions test-hooks

build: test
	go build \
		-ldflags "-X github.com/Skryensya/footprint/internal/app.Version=$(VERSION)" \
		-o fp \
		./cmd/fp

test:
	go test -count=1 ./...

test-actions:
	go test -count=1 ./internal/actions

test-hooks: build
	./scripts/test-hooks.sh
