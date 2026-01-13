VERSION := $(shell git describe --tags --dirty --always)

build:
	go build \
		-ldflags "-X github.com/Skryensya/footprint/internal/app.Version=$(VERSION)" \
		-o fp \
		./cmd/fp
