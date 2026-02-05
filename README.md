# forge-workhorse-groningen

> A Fulmen workhorse application template for robust, scalable Go backends

Named after the Groningen horse breed from the Netherlands, renowned for strength and toughness in heavy work—originally helping with canals and plowing in heavy wet soil. The binary is simply called `groningen`.

## Overview

`forge-workhorse-groningen` is a **Level 2 template** in the Fulmen ecosystem—a production-ready starter that provides:

- ✅ HTTP server with standard endpoints (`/health`, `/version`, `/metrics`)
- ✅ Control plane / data plane split with dedicated operational server
- ✅ Optional data plane auth (bearer token, basic auth, route policy)
- ✅ CLI with required subcommands (serve, version, health, envinfo, doctor)
- ✅ Structured logging with progressive profiles (via gofulmen)
- ✅ Three-layer configuration with canonical/alias env var ergonomics
- ✅ Graceful shutdown and signal handling
- ✅ Observability and telemetry built-in
- ✅ CRDL philosophy: Clone → Degit → Refit → Launch

## Fulmen Ecosystem Layers

```
Level 3: Your Application ← You are here after refitting
Level 2: forge-workhorse-groningen ← We are here (template)
Level 1: gofulmen + goneat (helpers + tooling)
Level 0: Crucible (SSOT - schemas, standards, docs)
```

## Quick Start

### Prerequisites

- Go 1.25+ ([install](https://go.dev/doc/install))
- golangci-lint ([install](https://golangci-lint.run/welcome/install/))

### Bootstrap

```bash
# Clone template
git clone https://github.com/fulmenhq/forge-workhorse-groningen.git my-app
cd my-app

# Install dependencies
# - Requires 'sfetch' (trust anchor); bootstrap prints install instructions if missing
make bootstrap

# Run server
make run
```

The server will start at `http://localhost:8080` with:

- Health checks: `http://localhost:8080/health/*` (live, ready, startup)
- Version info: `http://localhost:8080/version`
- Metrics: `http://localhost:8080/metrics` (JSON) and `http://localhost:9090/metrics` (Prometheus)

## Architecture

### Directory Structure

```
forge-workhorse-groningen/
├── cmd/
│   └── groningen/              # Entry point
│       └── main.go             # Minimal main (version injection)
├── internal/
│   ├── cmd/                    # Cobra commands
│   │   ├── root.go             # Root command + global flags
│   │   ├── serve.go            # HTTP server command
│   │   ├── version.go          # Version command
│   │   ├── health.go           # Health self-check
│   │   ├── envinfo.go          # Environment info
│   │   └── doctor.go           # Diagnostics
│   ├── server/                 # HTTP server implementation
│   │   ├── server.go
│   │   ├── handlers/           # Health, version, metrics
│   │   ├── middleware/         # Logging, correlation IDs
│   │   ├── auth/              # Data plane auth middleware + route policy
│   │   └── control/           # Control plane server + handlers
│   ├── core/                   # Business logic (your code here)
│   ├── config/                 # Config management + env var mappings
│   └── observability/          # Logging, metrics setup
├── config/
│   └── groningen/
│       └── v1.0.0/
│           └── groningen-defaults.yaml  # Template defaults (Layer 1)
├── schemas/
│   └── groningen/
│       └── v1.0.0/
│           └── config.schema.json       # Config validation schema
├── docs/
│   ├── groningen-overview.md     # Template architecture and components
│   ├── schema-validation.md      # JSONSchema multi-draft guide
│   ├── metrics.md                # Telemetry/metrics notes
│   ├── decisions/                # Architecture Decision Records (ADRs)
│   ├── releases/                 # Release notes by version
│   └── development/              # Development guides
│       ├── README.md             # Development handbook and workflows
│       ├── ci.md                 # CI/CD details
│       └── fulmen_cdrl_guide.md  # How to refit this template
├── .env.example                # Standard env vars (copy to .env)
├── Makefile                    # Development targets
└── go.mod                      # Dependencies
```

### Dependencies

- **gofulmen** - Fulmen helper library (config, logging, telemetry, schema validation, pathfinder)
- **goneat** - Tooling for formatting, hooks, and assessment
- **sfetch** - Trust anchor for bootstrapping tools
- **cobra** - CLI framework (Fulmen standard for Go)
- **chi** - HTTP router (lightweight, idiomatic)

## CLI Commands

```bash
# Server management
groningen serve                 # Start data plane + control plane
groningen serve --port 9000     # Custom data plane port

# Information commands
groningen version               # Basic version
groningen version --extended    # Full version + SSOT info
groningen health                # Self-check
groningen envinfo               # Config, env var mapping table, SSOT

# Diagnostics
groningen doctor                # Run checks: env conflicts, auth, control plane
```

## Configuration

Groningen uses **gofulmen/config** for canonical three-layer configuration with schema validation.

### Config Loader

Configuration is managed by `internal/config/loader.go`, which implements:

- **Typed config structs** (no runtime string lookups)
- **Schema validation** via JSON Schema
- **Absolute path resolution** (works from any directory, including tests)
- **Graceful reload** on SIGHUP with logger reconfiguration

### Three-Layer Config Pattern

1. **Layer 1 (Template Defaults)**: `config/groningen/v1.0.0/groningen-defaults.yaml`
   - Shipped with the template
   - Provides sensible defaults for all configuration options
   - Validated against `schemas/groningen/v1.0.0/config.schema.json`

2. **Layer 2 (User Overrides)**: `~/.config/<vendor>/<binary-name>.yaml`
   - Discovered via app identity (`.fulmen/app.yaml`)
   - Merged on top of template defaults
   - Optional (falls back to defaults if not present)

3. **Layer 3 (Runtime Overrides)**: Environment variables and CLI flags
   - Highest priority
   - Environment variables use prefix from app identity (default: `GRONINGEN_`)
   - CLI flags override everything

**Priority**: CLI flags > Environment variables > User config > Template defaults

### Schema Validation

Configuration is validated against the JSON Schema at:

```
schemas/groningen/v1.0.0/config.schema.json
```

Validation happens on load and reload. Invalid configuration prevents application startup or reload (falls back to previous valid config).

### Environment Variables

All env vars use the prefix from app identity (default: `GRONINGEN_`). Both canonical (nested) and alias (short) forms are supported:

| Canonical (nested)                     | Alias (short)                        | Config Path                |
| -------------------------------------- | ------------------------------------ | -------------------------- |
| `GRONINGEN_SERVER_HOST`                | `GRONINGEN_HOST`                     | `server.host`              |
| `GRONINGEN_SERVER_PORT`                | `GRONINGEN_PORT`                     | `server.port`              |
| `GRONINGEN_LOGGING_LEVEL`              | `GRONINGEN_LOG_LEVEL`                | `logging.level`            |
| `GRONINGEN_LOGGING_PROFILE`            | `GRONINGEN_LOG_PROFILE`              | `logging.profile`          |
| `GRONINGEN_CONTROL_PLANE_PORT`         | `GRONINGEN_CONTROLPLANE_PORT`        | `controlPlane.port`        |
| `GRONINGEN_CONTROL_PLANE_BEARER_TOKEN` | `GRONINGEN_CONTROLPLANE_BEARERTOKEN` | `controlPlane.bearerToken` |
| `GRONINGEN_DATA_PLANE_AUTH_ENABLED`    | `GRONINGEN_DATAPLANEAUTH_ENABLED`    | `dataPlaneAuth.enabled`    |

**Conflict detection**: If both canonical and alias are set with different values, a warning is logged and the alias takes precedence. Run `groningen doctor` to detect conflicts, or `groningen envinfo` for a full mapping table.

**Sensitive value masking**: Values for env vars containing `TOKEN`, `SECRET`, `PASSWORD`, or `KEY` are displayed as `[set]` in diagnostic output.

Copy `.env.example` to `.env` and customize for local development. See `.env.example` for the full list with both canonical and alias forms.

## Development

### Make Targets

```bash
make help          # Show all targets
make bootstrap     # Install dependencies (first-time setup)
make build         # Build binary
make build-all     # Build for multiple platforms
make run           # Run in development mode
make test          # Run tests
make test-cov      # Run tests with coverage
make lint          # Run linting
make fmt           # Format code
make clean         # Clean build artifacts
make check-all     # Run lint + test
make version       # Print current version
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-cov

# Run specific package
go test ./internal/config/...
```

### Linting

```bash
# Run all linters
make lint

# Auto-fix issues
golangci-lint run --fix
```

## CRDL: Refit This Template

To create your own application from this template:

1. **Clone** the template:

   ```bash
   git clone https://github.com/fulmenhq/forge-workhorse-groningen.git my-app
   cd my-app
   ```

2. **Degit** (remove template git history):

   ```bash
   rm -rf .git
   git init
   ```

3. **Refit** (customize for your app):

   **Step 1: Update App Identity** (`.fulmen/app.yaml`)

   ```yaml
   vendor: mycompany # Your organization
   binary_name: myapi # Your application name
   service_type: workhorse # Keep this for workhorse templates
   env_prefix: MYAPI # Your env var prefix (uppercase)
   config_name: myapi # Your config file name
   ```

   **Step 2: Update Module Path**
   - Update `go.mod`: Change module path to your repository
   - Example: `module github.com/mycompany/myapi`

   **Step 3: Update Environment Variables**
   - Customize `.env.example` with your variables
   - No need to rename prefixes - app identity handles this!

   **Step 4: Update Config and Schema Files**
   - See [CDRL Guide](docs/development/fulmen_cdrl_guide.md) for detailed config/schema renaming instructions

   **Step 5: Customize Application**
   - Replace placeholder business logic in `internal/core/`
   - Update `README.md`, `LICENSE`, etc.
   - Update CLI command descriptions in `internal/cmd/`

4. **Launch**:
   ```bash
   make bootstrap
   make run
   ```

**Key Benefit**: With App Identity integration, you only need to update `.fulmen/app.yaml` and the codebase automatically uses your new identity for env vars, config paths, and telemetry namespaces!

## Observability

### Logging

Uses gofulmen's progressive logging profiles:

- **SIMPLE**: Console output for CLI (default for commands)
- **STRUCTURED**: JSON output with correlation IDs (default for server)
- **ENTERPRISE**: Full envelope with middleware, throttling, policy enforcement

Configure via:

- Config file: `logging.profile: "structured"`
- Environment: `GRONINGEN_LOG_LEVEL=debug`
- CLI flag: `--verbose`

### Metrics

Prometheus metrics exposed at `/metrics` (JSON endpoint on port 8080, Prometheus format on port 9090):

- `http_requests_total` - Total HTTP requests by method/path/status
- `http_request_duration_ms` - Request latency histogram
- Request ID correlation for tracing
- Standard Go runtime metrics (goroutines, memory, etc.)

**Request ID Correlation**: Every request gets a unique X-Request-ID header for tracing and debugging.

### Tracing

Optional OpenTelemetry integration (TBD).

## Production Reliability

### Graceful Shutdown

Groningen implements production-grade signal handling with graceful shutdown:

```bash
# Start server
groningen serve

# Graceful shutdown (SIGINT/SIGTERM)
# Ctrl+C or kill <pid>
# - Stops accepting new requests
# - Completes in-flight requests
# - Closes database connections
# - Flushes logs and metrics
# - Clean exit

# Force quit (double-tap)
# Press Ctrl+C twice within 2 seconds
# Immediate exit if shutdown hangs
```

**Shutdown Sequence** (LIFO order):

1. Stop accepting new connections on both data plane and control plane
2. Shutdown HTTP servers (wait for in-flight requests)
3. Flush logger (ensure all logs written)
4. Exit cleanly

### Config Reload

Send SIGHUP to reload configuration without restart:

```bash
# Send SIGHUP signal
kill -HUP $(pgrep groningen)

# Config reload attempts to re-read config file and apply changes
# Some changes may still require restart (e.g., port changes)
```

### Control Plane

A dedicated operational server runs on a separate port from the data plane, enforcing security isolation:

```
Data Plane (0.0.0.0:8080)    Control Plane (127.0.0.1:9091)
├── /health/*                 ├── /control/           (discovery)
├── /version                  ├── /control/signal     (signal injection)
├── /metrics                  └── /control/config/reload
└── /* (your routes)
```

```yaml
controlPlane:
  enabled: true # Enabled by default
  host: 127.0.0.1 # Loopback by default (safe)
  port: 9091
  basePath: /control
  bearerToken: "" # Required if host is not loopback
```

```bash
groningen serve

# Trigger config reload
curl -X POST http://127.0.0.1:9091/control/config/reload \
  -H "Authorization: Bearer your-secret-token"

# Inject a signal (allowed: SIGHUP, SIGTERM)
curl -X POST http://127.0.0.1:9091/control/signal \
  -H "Authorization: Bearer your-secret-token" \
  -d '{"signal": "SIGHUP"}'
```

**Security**: Control plane binds to loopback by default. If you bind to a non-loopback interface, a bearer token MUST be configured. See [ADR-0001](docs/decisions/ADR-0001-control-plane-data-plane-split.md) and [ADR-0002](docs/decisions/ADR-0002-control-plane-auth-localhost.md) for design rationale.

### Data Plane Auth (Optional)

Optional starter authentication for the primary API surface. Disabled by default. Designed as an on-ramp — CDRL users are expected to replace with their preferred auth strategy (OAuth, JWT, etc.).

**Auth modes**: `disabled` (default), `bearerToken`, `basicAuth`

**Route policy** with longest-prefix matching:

| Category      | Behavior                                             |
| ------------- | ---------------------------------------------------- |
| `deny`        | Return 404 (hide route existence)                    |
| `public`      | No auth required                                     |
| `conditional` | Auth optional; handlers get auth context if provided |
| `protected`   | 401 if not authenticated                             |

```yaml
dataPlaneAuth:
  enabled: false # Disabled by default
  mode: disabled # disabled | bearerToken | basicAuth
  bearerToken: ""
  basicAuth:
    username: ""
    password: ""
  routePolicy:
    - prefix: /health
      category: public
    - prefix: /version
      category: public
    - prefix: /metrics
      category: conditional
    - prefix: /
      category: protected # Catch-all
```

Handlers can inspect auth state via `auth.Get(r.Context())` for conditional logic. See [ADR-0003](docs/decisions/ADR-0003-data-plane-auth-route-policy.md) for design rationale.

### Exit Codes

Groningen uses standardized exit codes from the Foundry catalog for operational clarity and better shell scripting support:

| Code | Name          | When                                                               |
| ---- | ------------- | ------------------------------------------------------------------ |
| 0    | Success       | Command completed successfully                                     |
| 1    | Failure       | Generic failure (default for unspecified errors)                   |
| 30   | ConfigInvalid | Configuration file is invalid or logger initialization failed      |
| 50   | FileNotFound  | Required file not found (e.g., `.fulmen/app.yaml`, home directory) |

**Usage in Shell Scripts:**

```bash
# Check exit codes for automation
groningen health
if [ $? -eq 0 ]; then
    echo "Service is healthy"
elif [ $? -eq 63 ]; then
    echo "Service unavailable"
fi

# Handle specific failures
groningen serve
exit_code=$?
case $exit_code in
    0)
        echo "Server stopped cleanly"
        ;;
    30)
        echo "Configuration error - check config file"
        ;;
    50)
        echo "Missing required file - check .fulmen/app.yaml"
        ;;
    *)
        echo "Server error (exit code: $exit_code)"
        ;;
esac
```

**Exit Code Metadata:**

All exit codes include metadata in error logs (code, name, description, category) to help with troubleshooting. When a fatal error occurs, you'll see:

```
FATAL: Failed to load app identity from .fulmen/app.yaml: file not found
Exit Code: 50 (FileNotFound) - Required file not found
```

**Future Exit Codes:**

As additional features are added, more semantic exit codes may be introduced:

- `40` (InvalidArgument) - Invalid command-line arguments
- `51` (PermissionDenied) - Permission denied errors
- `62` (NetworkError) - Network connectivity issues
- `63` (ServiceUnavailable) - Service or dependency unavailable

## Standard Endpoints

### Health Checks

- `GET /health` – Aggregate of all registered checks with semantic status (`healthy`, `degraded`, `unhealthy`). Returns `503` when any dependency is unhealthy.
- `GET /health/live` – Liveness probe with fast timeout to ensure the process is still running.
- `GET /health/ready` – Readiness probe that ensures telemetry, signal handlers, and identity have finished initializing.
- `GET /health/startup` – Confirms initialization completed; useful for Kubernetes startup probes.

Each response includes version metadata, RFC3339 timestamps, and per-check statuses to simplify debugging.

### Version Information

- `GET /version` – Returns app identity (binary name, semantic version), git commit, build date, Go runtime info, and the embedded gofulmen/Crucible dependency versions pulled directly from the SSOT.

### Metrics

- `GET /metrics` – Exposes full Prometheus/OpenMetrics output in `text/plain; version=0.0.4` format by proxying the internal gofulmen exporter. Scrape this endpoint from the main HTTP port; it automatically respects the configured metrics port/namespace.

### Standardized Errors

All non-2xx responses use a consistent JSON envelope:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "The requested resource was not found"
  }
}
```

These helpers are wired into the chi router for 404/405 cases and can be reused by downstream handlers for custom errors.

## Current Status

✅ **v0.1.10** - Production-ready workhorse template

- [x] Project structure and dependencies
- [x] Root command with global flags
- [x] Configuration management (gofulmen/config + three-layer pattern)
- [x] Canonical/alias env var ergonomics with conflict detection
- [x] Serve command (data plane + control plane)
- [x] Control plane with signal injection and config reload
- [x] Optional data plane auth (bearer token, basic auth, route policy)
- [x] Health endpoints (live, ready, startup)
- [x] Version endpoint (full build info)
- [x] Metrics endpoint with Prometheus
- [x] Graceful shutdown with signal handling (both planes)
- [x] Version command (basic + extended)
- [x] Health command (CLI self-check)
- [x] Envinfo command (env var mapping table, sensitive masking)
- [x] Doctor command (env conflicts, auth/control plane validation)
- [x] App Identity integration
- [x] Exit codes with semantic meanings
- [x] Request ID correlation middleware
- [x] Standardized error handling
- [x] Config reload via SIGHUP (signal or control plane)
- [x] Full JSONSchema multi-draft support (Draft-04 through 2020-12)
- [x] Architecture Decision Records
- [x] Integration with gofulmen logging
- [x] Integration with gofulmen telemetry
- [x] Comprehensive tests
- [x] Documentation

## Contributing

See [MAINTAINERS.md](MAINTAINERS.md) for governance and project team information.

## Resources

### FulmenHQ Ecosystem

- [Crucible](https://github.com/fulmenhq/crucible) - SSOT for schemas, standards, docs
- [Gofulmen](https://github.com/fulmenhq/gofulmen) - Go helper library
- [Goneat](https://github.com/fulmenhq/goneat) - DX CLI tool

### Documentation

- [Template Overview](docs/groningen-overview.md) - Comprehensive guide to template architecture and components
- [Schema Validation Guide](docs/schema-validation.md) - JSONSchema multi-draft support (Draft-04 through 2020-12)
- [Architecture Decisions (ADRs)](docs/decisions/README.md) - Key design decisions and trade-offs for CDRL users
- [Developer Handbook](docs/development/README.md) - Development setup and workflows
- [Development Guides](docs/development/) - Focused guides for specific topics

### Standards Applied

Groningen implements standards from the Fulmen ecosystem including the Forge Workhorse Standard, Go coding standards, CLI structure patterns, and HTTP REST standards. See the [Crucible repository](https://github.com/fulmenhq/crucible) for complete standard specifications.

## License

Licensed under the MIT License. See [LICENSE](LICENSE) file for complete details.

**Trademarks**: "Fulmen" and "3 Leaps" are trademarks of 3 Leaps, LLC. While code is open source, please use distinct names for derivative works to prevent confusion. See LICENSE for full guidelines.

### OSS Policies (Organization-wide)

- Authoritative policies repository: https://github.com/3leaps/oss-policies/
- Code of Conduct: https://github.com/3leaps/oss-policies/blob/main/CODE_OF_CONDUCT.md
- Security Policy: https://github.com/3leaps/oss-policies/blob/main/SECURITY.md
- Contributing Guide: https://github.com/3leaps/oss-policies/blob/main/CONTRIBUTING.md

---

<div align="center">

⚡ **Strong. Reliable. Production-Ready.** ⚡

_Workhorse template for the FulmenHQ ecosystem_

<br><br>

**Built with ⚡ by the 3 Leaps team**
**Part of the [Fulmen Ecosystem](https://fulmenhq.dev) - Lightning-fast enterprise development**

**Level 2 Template** • **Production Ready** • **Batteries Included**

</div>
