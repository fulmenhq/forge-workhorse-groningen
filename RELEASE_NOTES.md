# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.2] - 2025-11-16

### Pathfinder Integration & Dependency Updates

**Release Type**: Patch Release (Dependency Update + Code Improvement)
**Status**: ✅ Released

#### Overview

This release replaces the manual repository root finding implementation with gofulmen/pathfinder's battle-tested `FindRepositoryRoot()` function, providing enhanced security, performance, and cross-language parity. Also updates gofulmen to v0.1.15 with logging redaction middleware and Crucible to v0.2.16.

#### Key Changes

**Pathfinder Integration** (`internal/config/loader.go`):

- **Replaced Manual Implementation**: Removed custom `findProjectRoot()` with pathfinder integration
- **Enhanced Security**: Home directory ceiling, symlink loop detection, multi-tenant isolation, container escape prevention
- **Performance**: <30µs for repository root discovery (vs manual upward traversal)
- **Code Reduction**: Removed 22 lines of duplicate code (36 lines → 14 lines)
- **Resolved TODO**: Addressed technical debt comment about using pathfinder when available

**Dependency Updates**:

- **gofulmen**: v0.1.14 → v0.1.15
  - New: Logging redaction middleware (PII/secrets filtering)
  - New: Pathfinder repository root discovery API
  - New: Schema validator fixes for subdirectory testing
  - Updated: Crucible v0.2.16 with logging middleware specs
- **crucible**: v0.2.14 → v0.2.16 (transitive via gofulmen)

#### Benefits

**Security Enhancements**:

- ✅ Home directory boundary prevents traversal above `$HOME`
- ✅ Symlink loop detection (TRAVERSAL_LOOP error with critical severity)
- ✅ Multi-tenant isolation (boundaries prevent cross-tenant data access)
- ✅ Container escape prevention
- ✅ Filesystem root detection (/, C:\, UNC paths)

**Performance**:

- ✅ <30µs for all operations (well under Crucible spec targets)
- ✅ 830x faster than spec for immediate match
- ✅ 367x-1,111x faster than spec for upward traversal

**Testing**:

- ✅ 36 pathfinder tests (9 basic + 17 security + 10 benchmarks)
- ✅ All existing tests pass unchanged
- ✅ Cross-language parity with tsfulmen v0.1.9

#### Files Modified

```
VERSION                          # 0.1.1 → 0.1.2
go.mod                           # gofulmen v0.1.15, crucible v0.2.16
go.sum                           # Updated checksums
internal/config/loader.go        # Pathfinder integration (-22 lines)
CHANGELOG.md                     # v0.1.2 entry
RELEASE_NOTES.md                 # This document
docs/releases/v0.1.2.md          # Archived release notes
```

#### Migration Notes

**No migration required** - fully backward compatible.

**For Developers**:

If you've customized `internal/config/loader.go`:

- The `findProjectRoot()` function signature is unchanged
- Behavior is identical (finds go.mod or .git)
- Enhanced with security boundaries and loop detection

#### Quality Metrics

- ✅ All tests passing (9 packages, 0 failures)
- ✅ Format: 0 issues (62 files checked)
- ✅ Lint: 0 issues (35 Go files checked)
- ✅ Security: 0 issues (govulncheck + gosec)
- ✅ Overall health: **100%**

#### Available gofulmen v0.1.15 Features

While this release only integrates pathfinder, gofulmen v0.1.15 also provides:

**Logging Redaction Middleware** (not yet integrated):

- Pattern-based PII/secrets filtering (API keys, passwords, SSNs, credit cards)
- Field-based redaction (password, token, secret, apiKey, etc.)
- Replacement modes: text `[REDACTED]` or hash (SHA-256 prefix)
- Helper functions: `BundleStructuredWithRedaction()`, `WithRedaction()`

These features are available for future integration if needed.

---

## [0.1.1] - 2025-11-15

### Documentation Corrections for Public Release

**Release Type**: Patch Release (Documentation Fixes)
**Status**: ✅ Released

#### Overview

This patch release corrects critical documentation inaccuracies identified during final review before the public repository release. No code changes - purely documentation corrections to ensure accuracy for public users.

#### Documentation Fixes

**README.md Corrections**:

- Remove outdated local gofulmen references (now using public v0.1.14)
- Update dependency version numbers (gofulmen v0.1.10 → v0.1.14, goneat v0.3.0+ → v0.3.2)
- Remove "WIP" markers from all completed features (serve, version, health, envinfo, doctor commands)
- Correct binary name throughout CLI examples (`workhorse` → `groningen`)
- Fix CDRL config/schema renaming instructions (now references comprehensive guide)
- Update configuration description (viper → gofulmen/config to reflect actual implementation)
- Fix MAINTAINERS.md link to point to local file
- Remove broken links to non-public Crucible documentation
- Update Standards section to reference public Crucible repository

**docs/groningen-overview.md Corrections**:

- Update Current Version: 0.1.0 → 0.1.1
- Update Gofulmen Version: 0.1.7 (local replace) → 0.1.14
- Update Crucible Version: 2025.10.5 → 0.2.14

#### Quality Assurance

- ✅ All documentation reviewed for accuracy
- ✅ Version references synchronized
- ✅ Links verified (internal and external)
- ✅ CLI examples tested with correct binary name
- ✅ CDRL workflow instructions accurate

#### Files Modified

```
VERSION                          # Bumped to 0.1.1
README.md                        # 13 corrections
docs/groningen-overview.md       # 3 version updates
CHANGELOG.md                     # v0.1.1 entry added
RELEASE_NOTES.md                 # This document
docs/releases/v0.1.1.md          # Archived release notes
```

#### Migration Notes

No migration required - documentation-only release.

#### Known Issues

None - documentation is now accurate for public release.

---

## [0.1.0] - 2025-11-15

### App Identity Module Integration

**Release Type**: Major Feature Integration
**Status**: 🚧 In Development

#### Overview

This release integrates the App Identity module from gofulmen v0.1.9 to bring Groningen into compliance with the Fulmen Forge Workhorse Standard. This establishes the foundation for standardized configuration management, environment variables, and telemetry namespacing.

#### Features

**Unified Error & Telemetry Pipeline** (`internal/errors`, `internal/server`, `internal/cmd`):

- Standardized HTTP error handling via shared gofulmen envelopes, ensuring consistent JSON responses, correlation IDs, logging, and metrics across handlers, middleware, and CLI call sites.
- Context-aware `Wrap*` helpers propagate request IDs from Cobra commands and HTTP middleware, eliminating synthetic fallbacks in audit scenarios.
- Health probes now package per-check detail within structured envelopes while maintaining compatibility with telemetry and alerting pipelines.

**Portable Metrics Testing** (`internal/server/middleware/metrics_test.go`, `test/integration/metrics_test.go`):

- Switched unit and integration suites to gofulmen's in-memory collectors and explicit IPv4 listeners so `go test ./...` passes inside restricted sandboxes.
- Added permission-aware guards for httptest servers, providing friendly skips when environments forbid socket binding.

**gofulmen/config Integration** (`internal/config/`, `config/groningen/v1.0.0/`, `schemas/groningen/v1.0.0/`):

- **Layered Configuration**: Three-layer config pattern (defaults → user file → env vars) with schema validation via gofulmen/config
- **Absolute Path Resolution**: Repository root detection for config/schema loading from any working directory (tests, CLI, subdirectories)
- **CDRL-Friendly Structure**: Versioned config defaults (`config/groningen/v1.0.0/groningen-defaults.yaml`) and schemas (`schemas/groningen/v1.0.0/config.schema.json`)
- **Type-Safe Access**: Replaced all Viper calls with typed config structs and mapstructure tags
- **Config Reload**: SIGHUP handler with validation and logger reconfiguration
- **Backward Compatibility**: Old config paths still work (XDG migration path maintained)

**App Identity Integration** (`internal/cmd/root.go`, `internal/observability/`):

- **Identity Loading**: Load app metadata from `.fulmen/app.yaml` at startup
- **Config Path Derivation**: Use `identity.ConfigParams()` for XDG-compliant paths
- **Env Var Management**: Use `identity.EnvVar()` for consistent variable naming
- **Telemetry Namespace**: Use `identity.TelemetryNamespace()` for metrics and logging
- **Backward Compatibility**: Old config paths still work (XDG migration path)

**CDRL Workflow Enhancement**:

- **Single-File Identity**: Users only update `.fulmen/app.yaml` to refit template
- **Simplified Refit**: No need to search/replace env var prefixes across codebase
- **Documentation**: Updated README with clear CDRL instructions

**Files Added**:

```
.fulmen/app.yaml                                    # App identity definition
config/groningen/v1.0.0/groningen-defaults.yaml     # Layer 1 config defaults
schemas/groningen/v1.0.0/config.schema.json         # Config validation schema
internal/config/config.go                           # Typed config structs
internal/config/loader.go                           # Config loader with path resolution
internal/config/loader_test.go                      # Config unit tests
internal/errors/                                    # Error handling package
docs/architecture/decisions/ADR-0001-*.md           # Repository root detection ADR
```

**Files Modified**:

```
internal/cmd/root.go                 # Load identity, config integration
internal/cmd/serve.go                # Signal handlers, config reload, telemetry
internal/observability/logger.go     # Accept optional telemetry namespace
internal/observability/metrics.go    # Accept optional telemetry namespace
internal/server/                     # Error envelopes, standard endpoints
Makefile                             # Added precommit/prepush targets
README.md                            # Updated CDRL section with identity workflow
go.mod                               # Upgraded gofulmen to v0.1.14
```

#### Quality Assurance

- ✅ **All Tests Passing**: internal/observability test suite (100% pass rate)
- ✅ **Zero Lint Issues**: goneat assess reports 0 issues (Excellent health)
- ✅ **Code Formatted**: All files formatted with goneat (26 files)
- ✅ **Build Successful**: Binary builds without errors
- ✅ **Manual Testing**: `./bin/groningen version` works with identity

#### Dependencies

- **gofulmen**: v0.1.14 (upgraded from v0.1.7) - App Identity, Signal Handling, Exit Codes, Config, Telemetry modules
- **crucible**: v0.2.14 (auto-upgraded from v0.2.1, transitive via gofulmen) - Schemas, standards, and validation
- **goneat**: v0.3.2 - Formatting, assessment, and git hooks integration

#### Migration Notes for Template Users

**No migration required** for existing Groningen deployments - this is template infrastructure.

**For new CDRL users** (recommended workflow):

1. Clone template: `git clone https://github.com/fulmenhq/forge-workhorse-groningen.git myapp`
2. Update `.fulmen/app.yaml`:
   ```yaml
   vendor: mycompany
   binary_name: myapi
   env_prefix: MYAPI
   config_name: myapi
   ```
3. Update `go.mod` module path
4. Run `make build` - application automatically uses new identity

**Key Benefit**: Identity changes in `.fulmen/app.yaml` automatically propagate to:

- Environment variable prefix (`MYAPI_*`)
- Config file paths (`~/.config/mycompany/myapi.yaml`)
- Telemetry namespace (`mycompany.myapi`)
- Logger service name

#### Known Limitations

- Identity is static per process (no dynamic reloading)
- Config path backward compatibility maintained (old paths checked first)

#### Next Steps

- Signal Handling Module integration (graceful shutdown, config reload)
- Foundry Exit Codes integration (standardized exit codes)
- Comprehensive integration testing with all three modules

---

## [0.0.1] - 2025-10-28

### Initial Template Bootstrap

**Release Type**: Initial Release
**Status**: ✅ Completed

#### Overview

Initial bootstrap of forge-workhorse-groningen template with gofulmen integration, HTTP server, CLI framework, and observability foundation.

#### Features

**Core Template Structure**:

- **HTTP Server**: Chi router with /health, /version, /metrics endpoints
- **CLI Framework**: Cobra commands (serve, version, health, envinfo, doctor)
- **Configuration**: Viper-based three-layer config (defaults → file → env vars)
- **Logging**: Gofulmen logging with SIMPLE (CLI) and STRUCTURED (server) profiles
- **Metrics**: Prometheus metrics via gofulmen telemetry
- **Graceful Shutdown**: Basic SIGINT/SIGTERM handling with timeout

**Gofulmen Integration**:

- **Version**: gofulmen v0.1.7
- **Crucible**: v0.2.1 (embedded via gofulmen)
- **Modules Used**: logging, telemetry, config

**Build Tooling**:

- **Makefile**: Comprehensive targets (build, test, lint, fmt, run)
- **goneat**: DX tooling v0.3.2 for formatting and assessment
- **Go Version**: 1.25.1

#### Quality Metrics

- ✅ **Tests**: All passing (internal/observability)
- ✅ **Build**: Binary builds successfully
- ✅ **Lint**: Clean (goneat assess)
- ✅ **Documentation**: README with CDRL guide

#### Files Structure

```
forge-workhorse-groningen/
├── cmd/groningen/main.go               # Main entry point
├── internal/
│   ├── cmd/                            # Cobra commands
│   │   ├── root.go                     # Root command & config
│   │   ├── serve.go                    # HTTP server command
│   │   ├── version.go                  # Version command
│   │   ├── health.go                   # Health check command
│   │   ├── envinfo.go                  # Environment info command
│   │   └── doctor.go                   # Diagnostic command
│   ├── observability/                  # Logging & metrics
│   │   ├── logger.go                   # Gofulmen logger init
│   │   ├── metrics.go                  # Gofulmen metrics init
│   │   └── gofulmen_test.go            # Integration tests
│   └── server/                         # HTTP server
│       ├── server.go                   # Server setup
│       ├── routes.go                   # Route definitions
│       ├── handlers/                   # HTTP handlers
│       └── middleware/                 # HTTP middleware
├── .env.example                        # Environment variable template
├── Makefile                            # Build automation
├── README.md                           # Template documentation
└── go.mod                              # Go module definition
```

#### Known Issues

- Hardcoded GRONINGEN\_ prefix (resolved in v0.1.0 with App Identity)
- Basic signal handling (enhanced in upcoming signal handling integration)
- No standardized exit codes (added in upcoming exit codes integration)
