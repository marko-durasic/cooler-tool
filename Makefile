# Makefile for Cooler Tool
# Provides various build targets with different optimization levels

BINARY_NAME=cooler
MAIN_PATH=./cmd/cooler
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go build flags
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION)"
LDFLAGS_DEV=-ldflags="-X main.version=$(VERSION)"

.PHONY: all build build-dev build-release clean test lint help

# Default target - optimized release build
all: build

# Development build - includes debug symbols
build-dev:
	@echo "Building development version..."
	go build $(LDFLAGS_DEV) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Development build complete: $(BINARY_NAME)"

# Standard optimized build
build:
	@echo "Building optimized version..."
	go build $(LDFLAGS) -trimpath -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"
	@ls -lh $(BINARY_NAME) | awk '{print "Binary size:", $$5}'

# Release build with maximum optimizations
build-release:
	@echo "Building release version with maximum optimizations..."
	CGO_ENABLED=0 go build $(LDFLAGS) -trimpath -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Release build complete: $(BINARY_NAME)"
	@ls -lh $(BINARY_NAME) | awk '{print "Binary size:", $$5}'

# Cross-compile for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -trimpath -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -trimpath -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -trimpath -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -trimpath -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH)

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint the code
lint:
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "Available targets:"
	@echo "  make build        - Build optimized binary (default)"
	@echo "  make build-dev    - Build with debug symbols"
	@echo "  make build-release- Build with maximum optimizations"
	@echo "  make build-all    - Cross-compile for all platforms"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage- Run tests with coverage report"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make tidy         - Tidy dependencies"
	@echo "  make clean        - Remove build artifacts"
