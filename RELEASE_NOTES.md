# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.10] - 2026-01-17

### Dependency Updates for Clone Readiness (Patch)

**Release Type**: Patch Release (Dependency Freshness)
**Status**: 🚧 Prepared

#### Overview

This patch brings all FulmenHQ and third-party dependencies up to their latest stable versions in preparation for cloning groningen as a base for new applications. No breaking changes; purely dependency version bumps with verified quality gates.

#### Key Changes

- **gofulmen**: v0.3.0 → v0.3.1 (Crucible v0.4.2 → v0.4.3 transitively)
- **chi**: v5.2.3 → v5.2.4 (HTTP router)
- **cobra**: v1.10.1 → v1.10.2 (CLI framework)
- **zap**: v1.27.0 → v1.27.1 (structured logging)

#### Quality Gates

All quality gates verified:

- `make lint` - 0 issues, 100% health
- `make test` - All 8 packages passing
- `make build` - Clean build

#### Migration Notes

No migration required. Drop-in compatible for existing consumers.

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

---

## [0.1.8] - 2025-12-19

### Embedded App Identity for Standalone Binaries (Patch)

**Release Type**: Patch Release (Artifact Contract + CDRL Reliability)
**Status**: ✅ Released

#### Overview

This patch makes the template’s built artifacts self-identifying by embedding app identity at build time. Basic CLI commands (like `version` and `--help`) now work when the binary is executed outside a repo checkout (e.g. copied to `/tmp` or installed on another machine) without requiring `.fulmen/app.yaml` to exist on disk.

#### Key Changes

- **Embedded identity fallback**: Mirrors `.fulmen/app.yaml` into an embeddable in-module path and registers it via gofulmen’s `RegisterEmbeddedIdentityYAML` so identity resolution works anywhere.
- **Drift guardrails**: Added `sync-embedded-identity` and `verify-embedded-identity` targets and wired sync into `build`, `test`, and `release-build`.
- **Conformance test**: Added an integration test that builds and runs the binary from a temp directory to prevent regressions.
- **Dependencies**: gofulmen v0.1.24 (Crucible v0.2.25 transitively).

#### Migration Notes

No migration required for template consumers. CDRL consumers should continue editing `.fulmen/app.yaml` as the SSOT; the build tooling keeps the embedded mirror in sync.
