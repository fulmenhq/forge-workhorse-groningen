# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.10] - 2026-01-29

### Full JSONSchema Flavor Support (Patch)

**Release Type**: Patch Release (Schema Validation + CDRL Capability)
**Status**: 🚧 Prepared

#### Overview

This release expands schema validation capabilities for CDRL users by exposing full JSONSchema draft support. The gofulmen/crucible dependency bump brings validation for all major JSONSchema drafts (Draft-04 through Draft 2020-12), enabling downstream applications to consume external schemas in any format without configuration.

#### Key Feature: Multi-Draft Schema Validation

CDRL users now inherit:

- **Auto-detected drafts**: The `$schema` field determines validation behavior automatically
- **All major drafts**: Draft-04, Draft-06, Draft-07, Draft 2019-09, Draft 2020-12
- **Meta-validation**: Validate your schemas are well-formed before deployment
- **Documented patterns**: See `docs/schema-validation.md` and `schemas/README.md`

#### Key Changes

- **gofulmen**: v0.3.0 → v0.3.2 (Crucible v0.4.2 → v0.4.9 transitively)
- **go**: 1.25.1 → 1.25.5
- **goneat**: v0.3.21 → v0.5.1
- **golang.org/x/text**: v0.30.0 → v0.33.0
- **Bootstrap**: Skip goneat install if already present (use FORCE=1 to reinstall)

#### New Files

- `docs/schema-validation.md` - User-facing schema validation guide
- `schemas/README.md` - Schema directory documentation
- `testdata/schemas/*.schema.json` - Test fixtures for all 5 supported drafts
- `internal/config/schema_flavors_test.go` - Integration tests demonstrating multi-flavor support

#### Quality Gates

All quality gates verified:

- `make lint` - 0 issues, 100% health
- `make test` - All packages passing (including new schema flavor tests)
- `make build` - Clean build

#### Migration Notes

No migration required. Drop-in compatible for existing consumers. CDRL users gain new schema capabilities automatically.

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
