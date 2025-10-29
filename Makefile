.PHONY: help bootstrap build build-all run test test-cov lint fmt clean check-all version version-bump install

# Binary and version information
BINARY_NAME := groningen
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)

help:  ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

bootstrap:  ## Install gofulmen and dependencies (first-time setup)
	@echo "→ Installing gofulmen from local path..."
	@cd ../gofulmen && go install ./...
	@echo "→ Downloading Go module dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Bootstrap complete! Run 'make run' to start the server."

install:  ## Install dependencies (alias for bootstrap)
	@$(MAKE) bootstrap

run:  ## Run server in development mode
	@go run ./cmd/$(BINARY_NAME) serve --verbose

build:  ## Build binary for current platform
	@echo "→ Building $(BINARY_NAME) v$(VERSION)..."
	@go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "✓ Binary built: bin/$(BINARY_NAME)"

build-all:  ## Build multi-platform binaries
	@echo "→ Building for multiple platforms..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME)
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/$(BINARY_NAME)
	@cd bin && sha256sum * > SHA256SUMS.txt 2>/dev/null || shasum -a 256 * > SHA256SUMS.txt
	@echo "✓ Multi-platform binaries built in bin/"

test:  ## Run tests
	@go test -v -race ./...

test-cov:  ## Run tests with coverage
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

lint:  ## Run linting
	@golangci-lint run

fmt:  ## Format code
	@gofmt -w .
	@goimports -w . 2>/dev/null || true

clean:  ## Clean build artifacts
	@rm -rf bin/ coverage.out coverage.html
	@echo "✓ Cleaned build artifacts"

check-all:  ## Run all checks (lint + test)
	@echo "→ Running lint..."
	@$(MAKE) lint
	@echo "→ Running tests..."
	@$(MAKE) test
	@echo "✓ All checks passed!"

version:  ## Print version
	@echo "$(VERSION)"

version-bump:  ## Bump version (CalVer - requires manual edit of VERSION file)
	@echo "Current version: $(VERSION)"
	@echo "To bump version, edit VERSION file manually with new CalVer version"
	@echo "Example: echo '2025.10.1' > VERSION"
