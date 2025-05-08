# Extract version and commit
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)
BUILD_TAGS := $(shell if echo $(VERSION) | grep -q "devnet"; then echo "devnet"; elif echo $(VERSION) | grep -q "testnet"; then echo "testnet"; else echo "mainnet"; fi)

# Build for Linux (AMD64)
build-linux:
	env GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS),linux,amd64" -o ./build/optio-linux ./cmd/optiod/main.go

# Build for macOS (Apple Silicon)
build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS),darwin,arm64" -o ./build/optio-darwin-arm64 ./cmd/optiod/main.go

# Build for macOS (Intel)
build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS),darwin,amd64" -o ./build/optio-darwin-amd64 ./cmd/optiod/main.go

# Build for both platforms
build: build-linux build-darwin-arm64 build-darwin-amd64
