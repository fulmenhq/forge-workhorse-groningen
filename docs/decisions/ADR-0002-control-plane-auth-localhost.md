# ADR-0002: Control Plane Auth and Loopback Defaults

**Status**: Accepted
**Date**: 2026-02-05
**Deciders**: @3leapsdave

## Context

Control-plane endpoints can trigger reload/shutdown behaviors and are therefore highly sensitive.
For template consumers, the safest default is "available locally for operators" without being
accidentally reachable from a network.

## Decision

- Default `controlPlane.host` to loopback (`127.0.0.1`).
- Require a bearer token when binding the control plane to a non-loopback interface.
- Apply bearer auth middleware to all control plane endpoints.

Implementation references:

- Validation: `internal/cmd/serve.go` (controlPlane validation)
- Auth middleware: `internal/server/control/middleware/auth.go`

## Consequences

- Safe-by-default for local development and typical operations.
- When exposing control plane beyond loopback, operators must explicitly configure auth.
- Token-based auth is intentionally minimal (CDRL users can replace with mTLS, OAuth, etc.).

## Alternatives Considered

- Always require a token even on loopback.
  - Rejected: increases friction for local development; still achievable by setting a token.
- mTLS-only.
  - Rejected: too heavy for a template default; recommended for production refits.
