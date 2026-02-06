# Changelog

All notable changes to this project will be documented in this file. Older entries are archived under `docs/releases/` once we ship tagged versions.

> **Maintenance**: Keep only the 10 most recent releases in reverse-chronological order. Purge older entries when adding new releases.

## [Unreleased]

## [0.1.11] - 2026-02-05

### Fixed

- **Versioning**: v0.1.10 was released without a VERSION bump. This release corrects the version number. All feature content is identical to v0.1.10.

## [0.1.10] - 2026-02-05

### Added

- **Control plane / data plane split**: Dedicated operational server on separate port (default `127.0.0.1:9091`). Endpoints: discovery, signal injection, config reload. Bearer token auth required for non-loopback. Code at `internal/server/control/`.
- **Starter auth framework**: Optional data plane auth with `bearerToken` and `basicAuth` modes. Route policy system (`deny`/`public`/`conditional`/`protected`). Disabled by default. Code at `internal/server/auth/`.
- **Env var canonical/alias ergonomics**: Both nested canonical (`GRONINGEN_SERVER_PORT`) and short alias (`GRONINGEN_PORT`) env var names. Conflict detection with warnings. Alias precedence. Code at `internal/config/envvars.go`.
- **Enhanced doctor/envinfo diagnostics**: Env var mapping table in `envinfo`, conflict detection in `doctor`, sensitive value masking for TOKEN/SECRET/PASSWORD/KEY patterns.
- **Full JSONSchema flavor support**: Schema validation supports all major drafts (Draft-04, Draft-06, Draft-07, Draft 2019-09, Draft 2020-12) with auto-detection from `$schema` field. Meta-validation, test fixtures, and documentation.
- **ADR documentation structure**: `docs/decisions/` with index, template, and 4 initial ADRs covering control plane split, control plane auth, data plane auth, and schema validation.

### Changed

- **Dependencies**: Upgraded gofulmen v0.3.1 → v0.3.3 (Crucible v0.4.3 → v0.4.9 transitively).
- **Dependencies**: Upgraded go 1.25.1 → 1.25.5, golang.org/x/text v0.30.0 → v0.33.0.
- **Dependencies**: Upgraded goneat v0.3.21 → v0.5.2.
- **Bootstrap**: Skip goneat install if already present (use FORCE=1 to reinstall).
- **README**: Added control plane, data plane auth, env var ergonomics, and ADR sections.

### Quality

- `make fmt`, `make lint`, `make test`, and `make build` verified for this release.

## [0.1.9] - 2025-12-20

### Added

- **CDRL guide clarity**: Documented template-only files to delete vs refit, plus common residue hotspots.

### Changed

- **Release target naming**: Standardized `release-*` prefixes for dist/release checksum and key verification targets (kept deprecated aliases).
- **Release checklist**: Default signing flow is now download CI artifacts, regenerate manifests, verify, sign, and upload provenance only.
- **Dependencies**: Upgraded gofulmen to v0.1.25 (Crucible v0.2.26 transitively).

## [0.1.8] - 2025-12-19

### Added

- **Embedded app identity**: Mirrors `.fulmen/app.yaml` into an embeddable path and registers it with gofulmen so distributed binaries can self-identify outside a repo checkout.
- **Drift guardrails**: Added `make sync-embedded-identity` and `make verify-embedded-identity` and wired sync into `build`, `test`, and `release-build`.
- **Standalone acceptance test**: Builds the binary, copies it into a temp directory, and verifies `version`/`--help` work without `.fulmen/app.yaml` present.

### Changed

- **Dependencies**: Upgraded gofulmen to v0.1.24 (Crucible v0.2.25 transitively).

## [0.1.7] - 2025-12-18

### Added

- **Release provenance workflow**: Added minisign-first (primary) manifest signing plus optional PGP, with dual manifests (`SHA256SUMS`, `SHA512SUMS`) staged under `dist/release/`.
- **Release upload modes**: `make release-upload` now uploads provenance outputs only; `make release-upload-all` exists for fully manual artifact publishing.
- **Checksum verification**: Added `make verify-checksums` to confirm manifests match staged artifacts.

### Fixed

- **Release workflow**: CI release publishing now uploads `dist/release/*` (avoids `bin/` footguns and duplicate checksum uploads).

### Changed

- **CDRL guidance**: Documented refitting the signing env var prefix (`<APP>_…`) and clarified `env_prefix` should include the trailing underscore.
- **Release checklist**: Expanded signing section to call out prep vs signing steps and the provenance-only upload default.

## [0.1.6] - 2025-12-17

### Fixed

- **CDRL hardcoded residue**: Removed template-name defaults from CLI surfaces and `/version` fallback, made tests CDRL-safe, and updated Makefile help/SBOM output to use `$(BINARY_NAME)`.
- **Developer docs**: Updated Crucible-docs access guide to be template-neutral and aligned with current toolchain.

### Changed

- **Dependencies**: Upgraded gofulmen to v0.1.22 (Crucible v0.2.23 transitively).
- **Release workflow assets**: Release workflow now uploads `bin/*` rather than hardcoding the template binary name.

## [0.1.5] - 2025-12-16

### Added

- **Release signing helpers**: Manual download/sign/upload scripts and Make targets to support offline/controlled signing.

### Fixed

- **Release workflow gating**: Release workflow now runs only for `refs/tags/v*` and no longer fails on normal `main` pushes.

### Changed

- **Tooling**: Configured goneat v0.3.21 `lint.shell.shfmt.args` so shell formatting is deterministic across machines.

## [0.1.4] - 2025-12-15

### Changed

- **Dependencies**: Upgraded gofulmen to v0.1.21 (transitively pulls Crucible v0.2.21 via gofulmen).
- **CI root discovery**: CI now exports `FULMEN_WORKSPACE_ROOT` (GitHub workspace) and the config loader uses it as a boundary hint in CI.
- **Release automation**: Added a tag-triggered release workflow to publish build artifacts to GitHub Releases.

## [0.1.3] - 2025-12-01

### Changed

- **Dependencies**: Upgraded gofulmen v0.1.15 → v0.1.20 (transitively pulls Crucible v0.2.20 via gofulmen)
- **Version**: Bumped template version to 0.1.3

### Quality

- `make fmt`, `make lint`, `make test`, and `make build` verified for this release.

## [0.1.2] - 2025-11-16

### Changed

- **Repository Root Discovery**: Replaced manual `findProjectRoot()` with `gofulmen/pathfinder.FindRepositoryRoot()` for improved security and robustness
- **Dependencies**: Updated gofulmen v0.1.14 → v0.1.15, crucible v0.2.14 → v0.2.16

### Improved

- **Security**: Pathfinder provides home directory ceiling, symlink loop detection, and multi-tenant isolation
- **Performance**: Repository root discovery now <30µs (well under spec targets)
- **Code Quality**: Removed 22 lines of duplicate code, resolved TODO comment
