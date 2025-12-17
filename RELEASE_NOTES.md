# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

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

## [0.1.3] - 2025-12-01

### Dependency Refresh (Patch)

**Release Type**: Patch Release (Dependency Update)
**Status**: ✅ Released

#### Overview

This patch bumps gofulmen to v0.1.20 (bringing embedded Crucible v0.2.20 transitively) and advances the template version to 0.1.3. Quality gates (fmt, lint, test) and `make build` pass.

#### Key Changes

- **Dependencies**: gofulmen v0.1.20 (Crucible v0.2.20 transitively)
- **Version**: Updated `VERSION` to 0.1.3

#### Migration Notes

No code changes are required for template consumers; pull the new version and continue using `.fulmen/app.yaml` for identity.
