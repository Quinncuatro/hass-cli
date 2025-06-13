Project Structure Requirements
File Organization
Root Directory Structure
hass-cli/
├── cmd/hass/main.go           # Application entry point
├── internal/                  # Private application packages
├── pkg/                       # Public library code (if any)
├── test/                      # Integration tests
├── docs/                      # Documentation
├── .project-rules/            # Development rules (this dir)
├── specs/                     # Functional specifications
├── scripts/                   # Build and utility scripts
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── Makefile                   # Build automation
├── README.md                  # Project documentation
├── LICENSE                    # License file
└── .gitignore                 # Git ignore rules
Internal Package Structure
internal/
├── app/                       # Application coordination
│   ├── app.go                # Main application type
│   ├── config.go             # Application configuration
│   └── version.go            # Version information
├── cli/                       # CLI command handling
│   ├── root.go               # Root command
│   ├── lights.go             # Light control commands
│   ├── status.go             # Status commands
│   └── config.go             # Config commands
├── tui/                       # Terminal UI
│   ├── app.go                # TUI application
│   ├── components/           # UI components
│   └── styles/               # Styling definitions
├── client/                    # Home Assistant API client
│   ├── client.go             # HTTP client
│   ├── entities.go           # Entity operations
│   └── websocket.go          # WebSocket client (future)
├── config/                    # Configuration management
│   ├── config.go             # Config types and loading
│   ├── keyring.go            # Secure credential storage
│   └── migration.go          # Config migration
├── entity/                    # Entity resolution
│   ├── resolver.go           # Entity matching logic
│   ├── fuzzy.go              # Fuzzy matching algorithms
│   └── types.go              # Entity type definitions
└── cache/                     # Caching layer
    ├── cache.go              # Cache interface
    ├── memory.go             # In-memory implementation
    └── disk.go               # Disk-based implementation
Naming Conventions
File Names
Use lowercase with underscores: entity_resolver.go
Test files: entity_resolver_test.go
Integration tests: client_integration_test.go
Benchmark tests: fuzzy_benchmark_test.go
Package Names
Short, descriptive, lowercase
Match directory name exactly
No underscores or mixed case
Examples: client, config, entity
Go Module Name
go
module github.com/yourusername/hass-cli

go 1.21
Build Configuration
Makefile Structure
makefile
.PHONY: build test clean install lint fmt

# Build variables
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Default target
all: test build

# Build the application
build:
	go build $(LDFLAGS) -o bin/hass ./cmd/hass

# Run tests with coverage
test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Install the application
install:
	go install $(LDFLAGS) ./cmd/hass

# Lint the code
lint:
	golangci-lint run ./...

# Format the code
fmt:
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out

# Development setup
dev-setup:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
Version Management
Version Information
go
// internal/app/version.go
package app

import (
    "fmt"
    "runtime"
)

var (
    Version   = "dev"        // Set by build process
    Commit    = "unknown"    // Set by build process
    BuildTime = "unknown"    // Set by build process
)

func VersionInfo() string {
    return fmt.Sprintf("hass-cli %s (%s) built at %s with %s",
        Version, Commit, BuildTime, runtime.Version())
}
Build Script Integration
bash
#!/bin/bash
# scripts/build.sh
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build -ldflags "
    -X github.com/yourusername/hass-cli/internal/app.Version=$VERSION
    -X github.com/yourusername/hass-cli/internal/app.Commit=$COMMIT
    -X github.com/yourusername/hass-cli/internal/app.BuildTime=$BUILD_TIME
" -o bin/hass ./cmd/hass
Documentation Structure
README.md Requirements
Project Description: What the tool does
Installation Instructions: How to install/build
Quick Start: Basic usage examples
Configuration: Setup instructions
Commands: Command reference
Contributing: Development guidelines
License: License information
Code Documentation
Package-level documentation for all packages
Function documentation for public APIs
Example tests for complex functionality
Inline comments for non-obvious logic
Git Configuration
.gitignore Requirements
gitignore
# Binaries
bin/
hass
hass.exe

# Test coverage
coverage.out
coverage.html

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Config files (may contain secrets)
config.yaml
.hass-cli/

# Build artifacts
dist/
Git Hooks
bash
#!/bin/sh
# .git/hooks/pre-commit
# Run tests and linting before commit

make fmt
make lint
make test

if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
Dependency Management
Go Module Guidelines
Use semantic versioning for releases
Pin dependencies to specific versions
Regular dependency updates with testing
Minimal external dependencies
Approved Dependencies
CLI Framework: github.com/spf13/cobra
TUI Framework: github.com/charmbracelet/bubbletea
HTTP Client: Standard library net/http
Testing: github.com/stretchr/testify
Keyring: github.com/zalando/go-keyring
Configuration: gopkg.in/yaml.v3
Release Process
Release Checklist
Update version in appropriate files
Update CHANGELOG.md
Run full test suite
Build for multiple platforms
Create git tag
Upload release binaries
Update documentation
Cross-Platform Builds
bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/hass-linux-amd64 ./cmd/hass
GOOS=darwin GOARCH=amd64 go build -o bin/hass-darwin-amd64 ./cmd/hass
GOOS=windows GOARCH=amd64 go build -o bin/hass-windows-amd64.exe ./cmd/hass
Quality Assurance
Pre-commit Requirements
All tests must pass
Code must be formatted (gofmt)
Linting must pass (golangci-lint)
No security vulnerabilities (gosec)
CI/CD Pipeline
Automated testing on multiple Go versions
Cross-platform build validation
Security scanning
Code coverage reporting
Automated releases on tags
