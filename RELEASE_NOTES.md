# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.11] - 2026-02-05

### Versioning Correction

**Release Type**: Patch (Version Correction)
**Status**: ✅ Released

v0.1.10 was released without a VERSION bump. This release corrects the version number. All feature content is identical to v0.1.10.

Full details: [`docs/releases/v0.1.11.md`](docs/releases/v0.1.11.md)

---

## [0.1.10] - 2026-02-05

### Control Plane, Starter Auth, Env Var Ergonomics, Diagnostics, Schema Flavors, ADRs

**Release Type**: Minor Feature Release (6 feature briefs)
**Status**: ✅ Released

#### Overview

The largest feature release since v0.1.0. Adds control plane/data plane separation, starter authentication, ergonomic environment variables with conflict detection, enhanced diagnostics, full JSONSchema multi-draft support, and Architecture Decision Records.

- **Control plane split**: Dedicated operational server on `127.0.0.1:9091` (loopback by default). Discovery, signal injection, config reload endpoints. Bearer token auth required for non-loopback.
- **Starter auth framework**: Optional data plane auth with `bearerToken`/`basicAuth` modes and route policy (`deny`/`public`/`conditional`/`protected`). Disabled by default.
- **Env var ergonomics**: Both canonical (`GRONINGEN_SERVER_PORT`) and alias (`GRONINGEN_PORT`) env var names. Conflict detection with warnings. Alias takes precedence.
- **Enhanced diagnostics**: Env var mapping table in `envinfo`, conflict/auth checks in `doctor`, sensitive value masking.
- **Full JSONSchema flavor support**: Draft-04 through Draft 2020-12 with auto-detection from `$schema` field.
- **ADR documentation**: 4 Architecture Decision Records in `docs/decisions/`.

#### Key Changes

- **gofulmen**: v0.3.1 → v0.3.3 (Crucible v0.4.3 → v0.4.9 transitively)
- **go**: 1.25.1 → 1.25.5
- **goneat**: v0.3.21 → v0.5.2
- **golang.org/x/text**: v0.30.0 → v0.33.0
- **Bootstrap**: Skip goneat install if already present (use `FORCE=1` to reinstall)

#### New Packages

- `internal/server/control/` — Control plane server, routes, handlers, middleware
- `internal/server/auth/` — Data plane auth middleware, route policy
- `internal/config/envvars.go` — Env var mapping with canonical/alias support

#### Quality Gates

All quality gates verified:

- `make fmt` — Clean
- `make lint` — 0 issues
- `make test` — All packages passing
- `make build` — Clean build

#### Migration Notes

No breaking changes. Drop-in compatible. Control plane enabled by default on loopback; data plane auth disabled by default; existing env vars continue as aliases.

Full details: [`docs/releases/v0.1.10.md`](docs/releases/v0.1.10.md)

---

## [0.1.9] - 2025-12-20

### Release Workflow Naming + CDRL Clarity (Patch)

**Release Type**: Patch Release (Release UX + CDRL Reliability)
**Status**: 🚧 Prepared

#### Overview

This patch smooths the release signing workflow by standardizing `release-*` target naming and clarifying the default provenance flow: download CI-built artifacts, regenerate checksum manifests locally, sign manifests, and upload only provenance assets (signatures, keys, manifests, notes). It also improves the CDRL guide with clearer guidance on what template-only files are safe to delete vs refit.

#### Key Changes

- **Release target naming**: Added `release-checksums`, `release-verify-checksums`, and `release-verify-keys` and kept deprecated aliases for one cycle.
- **Release checklist defaults**: Recommended path is now `release-clean → release-download → release-checksums → release-verify-checksums → release-sign → release-export-keys → release-verify-keys → release-notes → release-upload`.
- **CDRL guide**: Clarified template-only deletions, emphasized env prefix residue scanning, and listed common residue hotspots.
- **Dependencies**: gofulmen v0.1.25 (Crucible v0.2.26 transitively).

#### Migration Notes

No migration required for template consumers.
