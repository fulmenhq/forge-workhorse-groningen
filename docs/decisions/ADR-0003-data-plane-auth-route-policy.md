# ADR-0003: Data Plane Auth Route Policy

**Status**: Accepted
**Date**: 2026-02-05
**Deciders**: @3leapsdave

## Context

Template consumers often need a minimal authentication layer early (dev/staging,
internal services) before integrating a full identity provider.
At the same time, a workhorse template must remain refittable: authentication choices
are application-specific.

## Decision

Provide an optional, starter data-plane auth middleware with a simple route policy.

- Disabled by default (`dataPlaneAuth.enabled: false`).
- Two starter modes:
  - `bearerToken`
  - `basicAuth`
- Prefix-based route categories (longest-prefix wins):
  - `deny` (respond with 404)
  - `public` (no auth)
  - `conditional` (auth optional; handlers may inspect auth context)
  - `protected` (auth required)

Implementation references:

- Policy/middleware: `internal/server/auth/middleware.go`, `internal/server/auth/policy.go`
- Data plane wiring: `internal/server/server.go`

## Consequences

- Provides a safe, minimal starting point for common cases.
- Encourages explicit policy definition over implicit "protect everything" behavior.
- Not a replacement for production-grade auth; intended to be replaced/extended.

## Alternatives Considered

- Hardcode protected/public endpoints.
  - Rejected: less flexible; harder for CDRL users to adjust.
- Full-featured RBAC/ABAC.
  - Rejected: too opinionated and too heavy for a template default.
