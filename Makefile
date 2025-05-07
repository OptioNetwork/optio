# Extract version and commit
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

# Build for Linux (AMD64)
build-linux:
	env GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/optio-linux ./cmd/optiod/main.go

# Build for macOS (Apple Silicon)
build-darwin:
	env GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./build/optio-darwin ./cmd/optiod/main.go

# Build for both platforms
build: build-linux build-darwin
