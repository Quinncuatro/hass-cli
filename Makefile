.PHONY: build test clean install lint fmt vet help

# Build variables
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X github.com/quinncuatro/hass-cli/internal/app.Version=$(VERSION) -X github.com/quinncuatro/hass-cli/internal/app.Commit=$(COMMIT) -X github.com/quinncuatro/hass-cli/internal/app.BuildTime=$(BUILD_TIME)"

# Default target
all: test build

# Build the application
build:
	@echo "Building hass-cli..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/hass ./cmd/hass

# Run tests with coverage
test:
	@echo "Running tests..."
	go test -race -coverprofile=coverage.out ./...
	@echo "Coverage report:"
	go tool cover -func=coverage.out

# Generate HTML coverage report
coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install the application
install:
	@echo "Installing hass-cli..."
	go install $(LDFLAGS) ./cmd/hass

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run 'make dev-setup' to install."; \
		exit 1; \
	fi

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ coverage.out coverage.html

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Cross-platform builds
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/hass-linux-amd64 ./cmd/hass
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/hass-darwin-amd64 ./cmd/hass
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/hass-darwin-arm64 ./cmd/hass
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/hass-windows-amd64.exe ./cmd/hass

# Check everything before commit
check: fmt vet lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  test       - Run tests with coverage"
	@echo "  coverage   - Generate HTML coverage report"
	@echo "  install    - Install the application"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  clean      - Clean build artifacts"
	@echo "  dev-setup  - Setup development environment"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  check      - Run all checks (fmt, vet, lint, test)"
	@echo "  help       - Show this help"