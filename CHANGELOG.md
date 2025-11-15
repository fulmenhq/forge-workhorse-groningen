# Changelog

All notable changes to this project will be documented in this file. Older entries are archived under `docs/releases/` once we ship tagged versions.

## [0.1.0] - Unreleased

### Added

- **App Identity Module**: Canonical binary name, environment prefix, config paths, and telemetry namespace derived from `.fulmen/app.yaml`
- **gofulmen/config Integration**: Layered configuration with schema validation, absolute path resolution, and CDRL-friendly structure (`config/groningen/v1.0.0/`, `schemas/groningen/v1.0.0/`)
- **Error Handling Pipeline**: Unified HTTP error handling with gofulmen error envelopes, shared responders, and context-aware correlation ID propagation
- **Telemetry Stack**: Production-ready metrics using gofulmen exporters, middleware, and sandbox-friendly fake collectors for automated tests
- **Signal Handling Module**: Graceful shutdown (SIGTERM/SIGINT), config reload (SIGHUP), and double-tap force quit with cross-platform support
- **Foundry Exit Codes**: Standardized exit codes (Success=0, ConfigInvalid=30, FileNotFound=50) with semantic meanings for operational clarity
- **Standard Endpoints**: Health checks (`/health`, `/health/live`, `/health/ready`, `/health/startup`), version info (`/version`), and metrics (`/metrics`)
- **goneat Git Hooks**: Pre-commit and pre-push validation with `make precommit` and `make prepush` targets

### Improved

- Health probe responses now emit structured JSON envelopes with per-check detail while preserving logging/metrics parity.
- Observability and metrics initialization expose the actual exporter port for reuse and diagnostics.

### Testing

- Comprehensive unit, integration, and middleware suites run cleanly in restricted sandboxes via in-memory collectors, IPv4 listeners, and guarded network setup.

### Dependencies

- **gofulmen**: Upgraded from v0.1.7 → v0.1.14 (App Identity, Signal Handling, Exit Codes, Config, Telemetry modules)
- **crucible**: Auto-upgraded v0.2.1 → v0.2.14 (transitive via gofulmen, provides schemas and standards)
- **goneat**: v0.3.2 for formatting and quality assessment

### Documentation

- Updated operational guidelines in `AGENTS.md` to reiterate the permanent gitignore policy for `.plans/` artifacts and reinforce quality gates
- Comprehensive CDRL guide in `docs/development/fulmen_cdrl_guide.md` with config/schema renaming instructions
- Template overview in `docs/groningen-overview.md` documenting all components and integration patterns
- ADR-0001 for repository root detection strategy
