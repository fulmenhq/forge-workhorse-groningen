# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.6] - 2025-12-17

### CDRL Hardening & Template Residue Cleanup (Patch)

**Release Type**: Patch Release (CDRL Reliability + Docs)
**Status**: ✅ Released

#### Overview

This patch reduces CDRL friction by removing template-name defaults from CLI surfaces and `/version` handler fallbacks, updating tests to be identity-driven, and making Makefile outputs more refit-friendly. It also refreshes the developer guide for accessing embedded Crucible docs via gofulmen.

#### Key Changes

- **CDRL hardening**:
  - CLI root defaults are template-neutral and overwritten by app identity.
  - `/version` fallback derives from the executable name rather than hardcoding `groningen`.
  - Unit tests avoid asserting template-specific identity values.
  - Makefile help banner and SBOM output use `$(BINARY_NAME)`.
- **Docs**: Updated Crucible-docs access guide wording and version prerequisites.
- **Release workflow**: Asset upload uses `bin/*` to avoid baking in a template binary name.
- **Dependencies**: gofulmen v0.1.22 (Crucible v0.2.23 transitively)
- **Version**: Updated `VERSION` to 0.1.6

#### Migration Notes

No migration required for template consumers.

---

## [0.1.5] - 2025-12-16

### CI + Release Workflow Completion (Patch)

**Release Type**: Patch Release (Dependency Update + CI/Release Tooling)
**Status**: ✅ Released

#### Overview

This patch completes the gofulmen bump to v0.1.21 and hardens CI repository-root discovery in container runners. It also finalizes the release workflow and introduces manual signing helpers, plus aligns shell linting with goneat v0.3.21 `shfmt` args override.

#### Key Changes

- **Dependencies**: gofulmen v0.1.21 (Crucible v0.2.21 transitively)
- **CI root boundary**: CI now exports `FULMEN_WORKSPACE_ROOT`, used as a boundary hint only under CI while still requiring repository markers (`go.mod`, `.git`)
- **Release automation**: Tag-triggered GitHub Release publishing plus manual download/sign/upload helpers
- **Tooling**: Configured `lint.shell.shfmt.args` (goneat v0.3.21) and reformatted scripts/hooks to match
- **Version**: Updated `VERSION` to 0.1.5

#### Migration Notes

No migration required for template consumers.

---

## [0.1.4] - 2025-12-15

### Dependency Refresh + Initial Release Workflow (Patch)

**Release Type**: Patch Release (Dependency Update + CI/Release Tooling)
**Status**: ⚠️ Tagged (superseded by 0.1.5)

#### Overview

This patch bumps gofulmen to v0.1.21 and adds the initial tag-triggered release workflow, plus CI repository-root boundary hints. This release is superseded by 0.1.5 due to follow-up release workflow config + signing improvements.

#### Key Changes

- **Dependencies**: gofulmen v0.1.21 (Crucible v0.2.21 transitively)
- **CI root boundary**: CI exports `FULMEN_WORKSPACE_ROOT`; loader uses it as a boundary hint only in CI
- **Release automation**: Initial tag-triggered GitHub Release workflow
- **Version**: Updated `VERSION` to 0.1.4

#### Migration Notes

No migration required for template consumers.

---
