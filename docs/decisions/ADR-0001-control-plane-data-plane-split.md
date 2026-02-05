# ADR-0001: Control Plane and Data Plane Split

**Status**: Accepted
**Date**: 2026-02-05
**Deciders**: @3leapsdave

## Context

Workhorse services need operational endpoints (reload/shutdown triggers, diagnostics)
without expanding the attack surface of the primary API.
Historically, templates often put these endpoints on the same listener as the data plane,
which makes "accidental exposure" more likely.

## Decision

Introduce a dedicated control-plane HTTP server, separate from the data-plane HTTP server.

- Data plane: main API endpoints on `server.host` / `server.port`.
- Control plane: operational endpoints on `controlPlane.host` / `controlPlane.port` under `controlPlane.basePath`.
- Default posture: control plane binds to loopback.

Implementation references:

- Data plane server: `internal/server/server.go`
- Control plane server: `internal/server/control/server.go`
- Control plane routes: `internal/server/control/routes.go`

## Consequences

- Reduced risk of exposing operational endpoints publicly.
- Cleaner separation of operational concerns (rate-limits/auth can differ).
- Slightly more complexity (two listeners; two shutdown paths).

## Alternatives Considered

- Single listener with strict routing/middleware separation.
  - Rejected: harder to guarantee operational endpoints are not reachable when misconfigured.
