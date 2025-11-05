# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.0] - 2025-11-05 (In Development)

### App Identity Module Integration

**Release Type**: Major Feature Integration
**Status**: 🚧 In Development

#### Overview

This release integrates the App Identity module from gofulmen v0.1.9 to bring Groningen into compliance with the Fulmen Forge Workhorse Standard. This establishes the foundation for standardized configuration management, environment variables, and telemetry namespacing.

#### Features

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

**Files Modified**:

```
.fulmen/app.yaml                     # NEW: App identity definition
internal/cmd/root.go                 # Load identity, use for config/env vars
internal/cmd/serve.go                # Pass namespace to logger/metrics
internal/observability/logger.go     # Accept optional telemetry namespace
internal/observability/metrics.go    # Accept optional telemetry namespace
README.md                            # Updated CDRL section with identity workflow
```

#### Quality Assurance

- ✅ **All Tests Passing**: internal/observability test suite (100% pass rate)
- ✅ **Zero Lint Issues**: goneat assess reports 0 issues (Excellent health)
- ✅ **Code Formatted**: All files formatted with goneat (26 files)
- ✅ **Build Successful**: Binary builds without errors
- ✅ **Manual Testing**: `./bin/groningen version` works with identity

#### Dependencies

- **gofulmen**: Linked locally (../gofulmen) for v0.1.9+ features
- **Crucible**: Auto-upgraded v0.2.1 → v0.2.6 (transitive, provides app identity schema)

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
