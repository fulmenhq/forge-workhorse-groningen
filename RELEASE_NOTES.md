# Release Notes

This document tracks release notes for forge-workhorse-groningen releases.

> **Convention**: Keep only the latest 3 releases here to prevent file bloat. Older releases are archived in `docs/releases/`.

## [0.1.3] - 2025-12-01

### Dependency Refresh (Patch)

**Release Type**: Patch Release (Dependency Update)
**Status**: ✅ Released

#### Overview

This patch bumps gofulmen to v0.1.20 (bringing embedded Crucible v0.2.20 transitively) and advances the template version to 0.1.3. Quality gates (fmt, lint, test) and `make build` pass.

#### Key Changes

- **Dependencies**: gofulmen v0.1.20 (Crucible v0.2.20 transitively)
- **Version**: Updated `VERSION` to 0.1.3

#### Benefits

- Picks up the latest gofulmen improvements while keeping the template aligned with the embedded Crucible snapshot.
- Confirms build and quality gates on the refreshed dependency set.

#### Files Modified

```
VERSION
go.mod
go.sum
CHANGELOG.md
RELEASE_NOTES.md
docs/releases/v0.1.3.md
```

#### Migration Notes

No code changes are required for template consumers; pull the new version and continue using `.fulmen/app.yaml` for identity.

---

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
