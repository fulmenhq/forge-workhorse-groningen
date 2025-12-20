# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

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

---

## [0.1.7] - 2025-12-18

### Release Signing Workflow Parity (Patch)

**Release Type**: Patch Release (Release Process Reliability)
**Status**: ✅ Released

#### Overview

This patch aligns Groningen’s manual release-signing workflow with Fulmen conventions: artifacts stage in `dist/release/`, checksum manifests are dual-generated (`SHA256SUMS` + `SHA512SUMS`), and signing is manifest-only (minisign primary, optional PGP). It also fixes the CI release publishing workflow to upload the staged `dist/release/*` set.

#### Key Changes

- **Release artifacts staging**: `make release-build` stages cross-platform binaries in `dist/release/` and generates checksum manifests.
- **Manifest-only signing**: `make release-sign` signs `SHA256SUMS`/`SHA512SUMS` with minisign and optionally PGP.
- **Trust anchors**: `make release-export-keys` exports minisign and PGP public keys into `dist/release/`.
- **Validation**: `make verify-checksums` verifies checksum manifests; `make verify-release-keys` verifies exported keys are public-only.
- **Upload safety**: `make release-upload` uploads provenance outputs only (manifests + signatures + public keys + release notes). Use `make release-upload-all` for fully manual artifact publishing.
- **CI release workflow**: tag-triggered release publishing now uploads `dist/release/*`.

#### Migration Notes

No migration required for template consumers.

---
