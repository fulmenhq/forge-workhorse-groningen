.PHONY: help bootstrap bootstrap-force tools sync dependencies verify-dependencies version-bump lint test build build-all clean fmt version check-all precommit prepush run install test-cov
.PHONY: version-set version-bump-major version-bump-minor version-bump-patch release-check release-prepare release-build

# Binary and version information
BINARY_NAME := groningen
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)

# Go related variables
GOCMD := go
GOTEST := $(GOCMD) test
GOFMT := $(GOCMD) fmt
GOMOD := $(GOCMD) mod

# Tooling
GONEAT_VERSION := v0.3.8
GONEAT_BIN := $(firstword $(wildcard ./bin/goneat) $(shell command -v goneat 2>/dev/null))

# Default target
all: fmt test

help:  ## Show this help message
	@echo 'Groningen - Available Make Targets'
	@echo ''
	@echo 'Required targets (Makefile Standard):'
	@echo '  help            - Show this help message'
	@echo '  bootstrap       - Install external tools (goneat) and dependencies'
	@echo '  bootstrap-force - Force reinstall external tools'
	@echo '  tools           - Verify external tools are available'
	@echo '  dependencies    - Generate SBOM for supply-chain security'
	@echo '  lint            - Run lint/format/style checks'
	@echo '  test            - Run all tests'
	@echo '  build           - Build distributable artifacts'
	@echo '  build-all       - Build multi-platform binaries'
	@echo '  clean           - Remove build artifacts and caches'
	@echo '  fmt             - Format code'
	@echo '  version         - Print current version'
	@echo '  version-set     - Set version to specific value'
	@echo '  version-bump-major - Bump major version'
	@echo '  version-bump-minor - Bump minor version'
	@echo '  version-bump-patch - Bump patch version'
	@echo '  release-check   - Run release checklist validation'
	@echo '  release-prepare - Prepare for release'
	@echo '  release-build   - Build release artifacts'
	@echo '  check-all       - Run all quality checks (fmt, lint, test)'
	@echo '  precommit       - Run pre-commit hooks (check-all)'
	@echo '  prepush         - Run pre-push hooks (check-all)'
	@echo ''
	@echo 'Additional targets:'
	@echo '  run             - Run server in development mode'
	@echo '  test-cov        - Run tests with coverage report'
	@echo ''

bootstrap:  ## Install external tools (goneat) and dependencies
	@echo "Installing external tools..."
	@if [ "$(FORCE)" = "1" ] || [ "$(FORCE)" = "true" ]; then \
		rm -f ./bin/goneat; \
	fi
	@if [ ! -x "./bin/goneat" ]; then \
		echo "→ Downloading goneat $(GONEAT_VERSION) to ./bin"; \
		./scripts/install-goneat.sh; \
	else \
		echo "→ goneat already available at ./bin/goneat"; \
	fi
	@echo "→ Downloading Go module dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Bootstrap completed. Use 'goneat' or add $$GOPATH/bin to PATH"

bootstrap-force:  ## Force reinstall external tools
	@$(MAKE) bootstrap FORCE=1

tools:  ## Verify external tools are available
	@echo "Verifying external tools..."
	@if command -v goneat >/dev/null 2>&1; then \
		echo "✅ goneat: $$(goneat --version 2>&1 | head -n1)"; \
	else \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@echo "✅ All tools verified"

sync:  ## Sync assets from Crucible SSOT (placeholder)
	@echo "⚠️  Groningen does not consume SSOT assets directly"
	@echo "✅ Sync target satisfied (no-op)"

dependencies:  ## Generate SBOM for supply-chain security
	@if ! command -v goneat >/dev/null 2>&1; then \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@echo "Generating Software Bill of Materials (SBOM)..."
	@goneat dependencies --sbom --sbom-output sbom/groningen.cdx.json
	@echo "✅ SBOM generated at sbom/groningen.cdx.json"

verify-dependencies:  ## Alias for dependencies (compatibility)
	@$(MAKE) dependencies

install:  ## Install dependencies (alias for bootstrap)
	@$(MAKE) bootstrap

run:  ## Run server in development mode
	@go run ./cmd/$(BINARY_NAME) serve --verbose

version-bump:  ## Bump version (usage: make version-bump TYPE=patch|minor|major|calver)
	@if ! command -v goneat >/dev/null 2>&1; then \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@if [ -z "$(TYPE)" ]; then \
		echo "❌ TYPE not specified. Usage: make version-bump TYPE=patch|minor|major|calver"; \
		exit 1; \
	fi
	@echo "Bumping version ($(TYPE))..."
	@goneat version bump $(TYPE)
	@echo "✅ Version bumped to $$(cat VERSION)"

version-set:  ## Set version to specific value (usage: make version-set VERSION=x.y.z)
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION not specified. Usage: make version-set VERSION=x.y.z"; \
		exit 1; \
	fi
	@echo "$(VERSION)" > VERSION
	@echo "✅ Version set to $(VERSION)"

version-bump-major:  ## Bump major version
	@$(MAKE) version-bump TYPE=major

version-bump-minor:  ## Bump minor version
	@$(MAKE) version-bump TYPE=minor

version-bump-patch:  ## Bump patch version
	@$(MAKE) version-bump TYPE=patch

release-check:  ## Run release checklist validation
	@echo "Running release checklist..."
	@$(MAKE) check-all
	@echo "✅ Release check passed"

release-prepare:  ## Prepare for release (tests, version bump)
	@echo "Preparing release..."
	@$(MAKE) check-all
	@echo "✅ Release preparation complete"

release-build: build-all  ## Build release artifacts (binaries + checksums)
	@echo "✅ Release build complete"

build:  ## Build binary for current platform
	@echo "→ Building $(BINARY_NAME) v$(VERSION)..."
	@go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "✓ Binary built: bin/$(BINARY_NAME)"

build-all:  ## Build multi-platform binaries and generate checksums
	@echo "→ Building for multiple platforms..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME)
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/$(BINARY_NAME)
	@GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/$(BINARY_NAME)
	@cd bin && (sha256sum * > SHA256SUMS.txt 2>/dev/null || shasum -a 256 * > SHA256SUMS.txt)
	@echo "✓ Multi-platform binaries built in bin/"

version:  ## Print current version
	@echo "$(VERSION)"

test:  ## Run all tests
	@echo "Running test suite..."
	$(GOTEST) ./... -v

test-cov:  ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

lint:  ## Run lint checks
	@if [ -z "$(GONEAT_BIN)" ]; then \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@echo "Running Go vet..."
	@$(GOCMD) vet ./...
	@echo "Running goneat assess..."
	@$(GONEAT_BIN) assess --categories lint
	@echo "✅ Lint checks passed"

fmt:  ## Format code with goneat
	@if [ -z "$(GONEAT_BIN)" ]; then \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@echo "Formatting with goneat..."
	@$(GONEAT_BIN) format
	@echo "✅ Formatting completed"

check-all: fmt lint test  ## Run all quality checks (ensures fmt, lint, test)
	@echo "✅ All quality checks passed"

precommit:  ## Run pre-commit hooks
	@if ! command -v goneat >/dev/null 2>&1; then \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@echo "Running pre-commit validation..."
	@goneat format
	@goneat assess --check --categories format,lint --fail-on critical
	@echo "✅ Pre-commit checks passed"

prepush:  ## Run pre-push hooks
	@if ! command -v goneat >/dev/null 2>&1; then \
		echo "❌ goneat not found. Run 'make bootstrap' first."; \
		exit 1; \
	fi
	@echo "Running pre-push validation..."
	@goneat format
	@goneat assess --check --categories format,lint,security --fail-on high
	@echo "✅ Pre-push checks passed"

clean:  ## Clean build artifacts and reports
	@echo "Cleaning artifacts..."
	rm -rf bin/ dist/ coverage.out coverage.html
	@echo "✅ Clean completed"
