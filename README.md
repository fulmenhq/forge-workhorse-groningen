# forge-workhorse-groningen

> A Fulmen workhorse application template for robust, scalable Go backends

Named after the Groningen horse breed from the Netherlands, renowned for strength and toughness in heavy work—originally helping with canals and plowing in heavy wet soil. The binary is simply called `groningen`.

## Overview

`forge-workhorse-groningen` is a **Level 2 template** in the Fulmen ecosystem—a production-ready starter that provides:

- ✅ HTTP server with standard endpoints (`/health`, `/version`, `/metrics`)
- ✅ CLI with required subcommands (serve, version, health, envinfo, doctor)
- ✅ Structured logging with progressive profiles (via gofulmen)
- ✅ Three-layer configuration management (Crucible → User → Runtime)
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

- Go 1.23+ ([install](https://go.dev/doc/install))
- golangci-lint ([install](https://golangci-lint.run/welcome/install/))
- Access to gofulmen (local at `../gofulmen`)

### Bootstrap

```bash
# Clone the template
git clone https://github.com/fulmenhq/forge-workhorse-groningen.git my-app
cd my-app

# Install dependencies (including gofulmen from local path)
make bootstrap

# Run the server
make run
```

The server will start at `http://localhost:8080` with:

- Health checks: `http://localhost:8080/health/*`
- Version info: `http://localhost:8080/version`
- Metrics: `http://localhost:9090/metrics`

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
│   │   ├── serve.go            # HTTP server command (WIP)
│   │   ├── version.go          # Version command (WIP)
│   │   ├── health.go           # Health self-check (WIP)
│   │   ├── envinfo.go          # Environment info (WIP)
│   │   └── doctor.go           # Diagnostics (WIP)
│   ├── server/                 # HTTP server implementation (WIP)
│   │   ├── server.go
│   │   ├── routes.go
│   │   ├── handlers/           # Health, version, metrics
│   │   └── middleware/         # Logging, correlation IDs
│   ├── core/                   # Business logic (your code here)
│   ├── config/                 # Config management
│   └── observability/          # Logging, metrics setup
├── config/
│   └── groningen.yaml          # App defaults (Layer 2)
├── docs/
│   ├── README.md
│   ├── DEVELOPMENT.md
│   └── development/
│       └── fulmen_cdrl_guide.md  # How to refit this template
├── .env.example                # Standard env vars (copy to .env)
├── Makefile                    # Development targets
└── go.mod                      # Dependencies
```

### Dependencies

- **gofulmen v0.1.5** - Fulmen helper library (config, logging, telemetry, etc.)
- **goneat v0.3.0** - Optional DX tooling
- **cobra** - CLI framework (Fulmen standard for Go)
- **viper** - Configuration management
- **chi** - HTTP router (lightweight, idiomatic)

## CLI Commands

```bash
# Server management
groningen serve                 # Start HTTP server
groningen serve --port 9000     # Custom port

# Information commands
groningen version               # Basic version
groningen version --extended    # Full version + SSOT info
groningen health                # Self-check
groningen envinfo               # Dump config/env/SSOT

# Diagnostics
groningen doctor                # Run checks, suggest fixes

# Configuration
groningen config show           # Display current config
groningen config validate       # Validate config file
```

## Configuration

### Three-Layer Config Pattern

1. **Layer 1 (Crucible)**: Default schemas and configs from SSOT (via gofulmen)
2. **Layer 2 (User)**: Your config file at `~/.config/workhorse-groningen/config.yaml`
3. **Layer 3 (Runtime)**: Environment variables and CLI flags

Priority: CLI flags > Environment variables > Config file > Crucible defaults

### Environment Variables

All env vars use the prefix `GRONINGEN_`:

```bash
GRONINGEN_PORT=8080
GRONINGEN_HOST=localhost
GRONINGEN_LOG_LEVEL=info
GRONINGEN_METRICS_PORT=9090
# ... see .env.example for full list
```

Copy `.env.example` to `.env` and customize for local development.

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
   - Rename `groningen` to `my-app` throughout (binary name, cmd directory)
   - Update `go.mod`: module path
   - Update env var prefix: `GRONINGEN_` → `MY_APP_`
   - Update `config/groningen.yaml` → `config/my-app.yaml`
   - Customize `.env.example` with your variables
   - Replace placeholder business logic in `internal/core/`
   - Update `README.md`, `LICENSE`, etc.

4. **Launch**:
   ```bash
   make bootstrap
   make run
   ```

See [docs/development/fulmen_cdrl_guide.md](docs/development/fulmen_cdrl_guide.md) for detailed refitting instructions.

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

Prometheus metrics exposed at `/metrics` (default port 9090):

- `http_requests_total` - Total HTTP requests by method/path/status
- `http_request_duration_seconds` - Request latency histogram
- Standard Go runtime metrics (goroutines, memory, etc.)

### Tracing

Optional OpenTelemetry integration (TBD).

## Standard Endpoints

### Health Checks

- `GET /health/live` - Liveness probe (200 if process is running)
- `GET /health/ready` - Readiness probe (200 if ready to serve traffic)
- `GET /health/startup` - Startup probe (200 when initialization complete)

### Version Information

- `GET /version` - Version info (app version, Crucible version, build info)

### Metrics

- `GET /metrics` - Prometheus metrics export

## Current Status

🚧 **Work in Progress** - Foundation complete, implementing commands and server

- [x] Project structure and dependencies
- [x] Root command with global flags
- [x] Configuration management (viper + three-layer pattern)
- [ ] Serve command (HTTP server with chi)
- [ ] Health endpoints
- [ ] Version endpoint
- [ ] Metrics endpoint with Prometheus
- [ ] Graceful shutdown
- [ ] Version command (basic + extended)
- [ ] Health command (CLI self-check)
- [ ] Envinfo command
- [ ] Doctor command
- [ ] Config subcommands
- [ ] Integration with gofulmen logging
- [ ] Integration with gofulmen telemetry
- [ ] Comprehensive tests
- [ ] Documentation

## Contributing

See [MAINTAINERS.md](../crucible/MAINTAINERS.md) for governance and [DEVELOPMENT.md](docs/DEVELOPMENT.md) for setup.

## Resources

### FulmenHQ Ecosystem

- [Crucible](https://github.com/fulmenhq/crucible) - SSOT for schemas, standards, docs
- [Gofulmen](https://github.com/fulmenhq/gofulmen) - Go helper library
- [Goneat](https://github.com/fulmenhq/goneat) - DX CLI tool
- [Technical Manifesto](../crucible/docs/architecture/fulmen-technical-manifesto.md)

### Standards Applied

- [Fulmen Workhorse Standard](../crucible/docs/architecture/fulmen-forge-workhorse-standard.md)
- [Go Coding Standards](../crucible/docs/standards/coding/go.md)
- [Go CLI (Cobra) Structure](../crucible/docs/standards/repository-structure/go/cli-cobra.md)
- [HTTP REST Standards](../crucible/docs/standards/api/http-rest-standards.md)

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
