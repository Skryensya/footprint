VERSION := $(shell git describe --tags --dirty --always)

build:
	go build \
	-ldflags "-X github.com/Skryensya/footprint/internal/app.Version=ba7024f-dirty" \
	-o fp \
	./cmd/fp