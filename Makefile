# Extract commit
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

# Build for Linux (AMD64)
build-linux:
	env GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=linux,amd64" -o ./build/optio-linux ./cmd/optiod/main.go

# Build for macOS (Apple Silicon)
build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=darwin,arm64" -o ./build/optio-darwin-arm64 ./cmd/optiod/main.go

# Build for macOS (Intel)
build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=darwin,amd64" -o ./build/optio-darwin-amd64 ./cmd/optiod/main.go

# Build for Windows (AMD64)
build-windows:
	env GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS) -X github.com/cosmos/cosmos-sdk/version.BuildTags=windows,amd64" -o ./build/optio-windows.exe ./cmd/optiod/main.go

# Build for all platforms
build: build-linux build-darwin-arm64 build-darwin-amd64 build-windows
