---
title: "Groningen Template Overview"
description: "Comprehensive overview of the Groningen workhorse application template"
author: "Forge Architect"
date: "2025-11-06"
last_updated: "2025-11-06"
status: "active"
tags: ["overview", "template", "workhorse", "go", "fulmen"]
---

# Groningen Template Overview

## Purpose & Scope

Groningen is a production-ready workhorse application template from the FulmenHQ ecosystem, providing enterprise-grade patterns for robust, scalable Go backends. It serves as the canonical implementation of the Fulmen Forge Workhorse Standard, demonstrating best practices for HTTP services, CLI tooling, observability, and operational excellence.

**Target Audience**: Go developers building production workloads, API services, background workers, and enterprise applications that need reliability, observability, and operational maturity from day one.

**Design Philosophy**: Progressive complexity with production-ready defaults. Simple deployments work out-of-the-box, while complex applications have access to full enterprise features including graceful shutdown, structured logging, metrics, and configuration management.

## Template Architecture

### Core Components

| Component           | Status              | Purpose                                                                   | Key Features |
| ------------------- | ------------------- | ------------------------------------------------------------------------- | ------------ |
| **HTTP Server**     | ✅ Production-ready | Chi router with standard endpoints, graceful shutdown, middleware support |
| **CLI Framework**   | ✅ Complete         | Cobra-based commands (serve, version, health, envinfo, doctor)            |
| **Configuration**   | ✅ Three-layer      | Defaults → file → environment variables with XDG compliance               |
| **Observability**   | ✅ Enterprise       | Structured logging, Prometheus metrics, health checks                     |
| **Signal Handling** | ✅ Cross-platform   | Graceful shutdown, config reload, double-tap force quit                   |
| **App Identity**    | ✅ CDRL-ready       | `.fulmen/app.yaml` for template customization                             |
| **Exit Codes**      | ✅ Standardized     | Foundry exit codes with semantic meaning                                  |

### Standard Endpoints

| Endpoint          | Method | Purpose                                          | Response                             |
| ----------------- | ------ | ------------------------------------------------ | ------------------------------------ |
| `/health`         | GET    | Aggregate application health across dependencies | JSON with status, timestamp, checks  |
| `/health/live`    | GET    | Fast liveness probe                              | JSON probe response (`200/503`)      |
| `/health/ready`   | GET    | Readiness probe w/ dependency checks             | JSON probe response (`200/503`)      |
| `/health/startup` | GET    | Startup probe to signal init completion          | JSON probe response (`200/503`)      |
| `/version`        | GET    | Version and identity information                 | JSON with app, SSOT, runtime details |
| `/metrics`        | GET    | Prometheus/OpenMetrics scrape endpoint           | OpenMetrics format for scraping      |

All non-2xx responses share the standardized error envelope described in §9 of the Workhorse standard, ensuring consistent JSON errors for orchestrators and operators.

The `/metrics` handler proxies directly to the gofulmen Prometheus exporter that runs on the configured metrics port, so platform teams only need to scrape the primary HTTP server to collect telemetry.

### CLI Commands

| Command   | Purpose              | Key Features                                      |
| --------- | -------------------- | ------------------------------------------------- |
| `serve`   | Start HTTP server    | Graceful shutdown, signal handling, config reload |
| `version` | Display version info | Basic and extended output with SSOT versions      |
| `health`  | Run health check     | Self-diagnosis with detailed status reporting     |
| `envinfo` | Show environment     | Comprehensive runtime and configuration info      |
| `doctor`  | System diagnostics   | Validates Go version, dependencies, config paths  |

## Integration Stack

### Fulmen Ecosystem Integration

| Module                   | Integration       | Benefits                                         |
| ------------------------ | ----------------- | ------------------------------------------------ |
| **gofulmen/appidentity** | ✅ Complete       | Dynamic app name, config paths, env var prefixes |
| **gofulmen/logging**     | ✅ Progressive    | SIMPLE → STRUCTURED → ENTERPRISE profiles        |
| **gofulmen/config**      | ✅ XDG-compliant  | Three-layer loading with validation              |
| **gofulmen/telemetry**   | ✅ Prometheus     | Counters, gauges, histograms with HTTP exporter  |
| **gofulmen/signals**     | ✅ Cross-platform | Graceful shutdown, SIGHUP reload, double-tap     |
| **gofulmen/foundry**     | ✅ Exit codes     | 54 standardized codes with metadata              |
| **gofulmen/crucible**    | ✅ Embedded       | Schema validation, standards access              |

### External Dependencies

| Dependency                   | Version | Purpose                              | Integration |
| ---------------------------- | ------- | ------------------------------------ | ----------- |
| **github.com/go-chi/chi/v5** | v5.2.3  | HTTP router with middleware support  |
| **github.com/spf13/cobra**   | v1.10.1 | CLI framework with command structure |
| **github.com/spf13/viper**   | v1.21.0 | Configuration management and binding |
| **go.uber.org/zap**          | v1.27.0 | High-performance structured logging  |

## Configuration Management

### Three-Layer Loading

1. **Defaults**: Built-in sensible defaults for development
2. **File**: Optional YAML config file with XDG-compliant paths
3. **Environment**: Environment variables override file settings

### Config File Discovery

```
$XDG_CONFIG_HOME/fulmen/groningen/config.yaml
~/.config/fulmen/groningen/config.yaml (fallback)
```

### Environment Variables

All environment variables use the app identity prefix:

```bash
# Server configuration
GRONINGEN_SERVER_HOST=localhost
GRONINGEN_SERVER_PORT=8080

# Logging configuration
GRONINGEN_LOG_LEVEL=info
GRONINGEN_LOG_PROFILE=structured

# Metrics configuration
GRONINGEN_METRICS_PORT=9090
GRONINGEN_METRICS_ENABLED=true
```

## Observability Stack

### Logging Profiles

| Profile        | Use Case                   | Output                                   | Features |
| -------------- | -------------------------- | ---------------------------------------- | -------- |
| **SIMPLE**     | CLI tools, development     | Console output, basic severity           |
| **STRUCTURED** | Production services        | JSON output, correlation IDs, file sinks |
| **ENTERPRISE** | Mission-critical workloads | Full envelope, middleware, throttling    |

### Metrics Collection

Built-in Prometheus metrics with standard exporters:

```go
// Counters for business metrics
requestCounter := telemetry.Counter("http_requests_total")
requestCounter.Inc()

// Gauges for system state
activeConnections := telemetry.Gauge("active_connections")
activeConnections.Set(42)

// Histograms for performance
requestDuration := telemetry.Histogram("http_request_duration_seconds")
requestDuration.Observe(0.123)
```

### Health Checks

Comprehensive health monitoring with:

- **Application Status**: Overall service health
- **Dependency Checks**: External service connectivity
- **Resource Metrics**: Memory, CPU, goroutine counts
- **Configuration Validation**: Required settings present

## Operational Patterns

### Graceful Shutdown

LIFO cleanup chain ensures proper resource release:

```go
signals.OnShutdown(func(ctx context.Context) error {
    logger.Info("Flushing metrics...")
    return telemetry.Flush()
})

signals.OnShutdown(func(ctx context.Context) error {
    logger.Info("Closing database...")
    return db.Close()
})
```

### Config Reload

SIGHUP handling for zero-downtime configuration updates:

```go
signals.OnReload(func(ctx context.Context) error {
    logger.Info("Reloading configuration...")
    newConfig, err := config.Load(configPath)
    if err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    config.Apply(newConfig)
    return nil
})
```

### Double-Tap Force Quit

Operator-friendly Ctrl+C handling with configurable window:

```go
signals.EnableDoubleTap(signals.DoubleTapConfig{
    Window:  2 * time.Second,
    Message: "Press Ctrl+C again to force quit",
})
```

## CDRL Workflow

### Template Customization

The template is designed for CDRL (Clone → Degit → Refit → Launch):

1. **Clone**: `git clone` the repository
2. **Degit**: Remove git history
3. **Refit**: Customize `.fulmen/app.yaml`:
   ```yaml
   app:
     vendor: yourcompany
     binary_name: yourapp
     env_prefix: YOURAPP_
     config_name: yourapp
   ```
4. **Launch**: Run your customized application

### Identity-Driven Customization

All application surfaces automatically adapt to your identity:

- **CLI Help**: Shows your app name and config paths
- **Logging**: Uses your app name as service identifier
- **Metrics**: Uses your telemetry namespace
- **Config**: Uses your vendor/config paths
- **Environment**: Uses your env var prefix

## Development Workflow

### Local Development

```bash
# Bootstrap development tools
make bootstrap

# Run tests
make test

# Build application
make build

# Run with hot reload
make run

# Lint and format
make lint
make fmt
```

### Quality Gates

```bash
# Full quality check
make check-all

# Individual checks
make test       # Run all tests
make lint       # Run Go vet and golangci-lint
make fmt        # Format code with goneat
make build      # Verify build succeeds
```

### Testing Strategy

- **Unit Tests**: Package-level testing with fixtures
- **Integration Tests**: End-to-end HTTP server testing
- **Observability Tests**: Logging and metrics validation
- **Signal Tests**: Graceful shutdown verification

## Production Deployment

### Container Support

Docker-ready with multi-stage builds:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/bin/groningen /usr/local/bin/
EXPOSE 8080 9090
CMD ["groningen", "serve"]
```

### Configuration Management

Production-ready configuration patterns:

- **XDG Compliance**: Standard config directory locations
- **Environment First**: Production settings via env vars
- **Validation**: Schema validation for all configuration
- **Reload Support**: Zero-downtime configuration updates

### Observability Integration

- **Prometheus**: Standard metrics endpoint on `/metrics`
- **Structured Logs**: JSON format with correlation IDs
- **Health Endpoints**: `/health` for load balancer checks
- **Graceful Shutdown**: Proper resource cleanup on termination

## Security Considerations

### Input Validation

- **Schema Validation**: All configuration validated against schemas
- **Request Validation**: HTTP inputs validated with middleware
- **Path Safety**: Filesystem operations protected against traversal

### Operational Security

- **Signal Handling**: Proper signal handling prevents resource leaks
- **Graceful Shutdown**: Ensures data integrity on termination
- **Config Reload**: Validates configuration before applying changes

## Performance Characteristics

### Resource Usage

- **Memory**: Efficient logging with zap backend
- **CPU**: Minimal overhead for observability
- **Network**: Efficient HTTP routing with chi
- **Disk**: Lazy config loading with validation caching

### Scalability Patterns

- **Horizontal**: Stateless design supports multiple instances
- **Vertical**: Configurable resource limits and pools
- **Observability**: Built-in metrics for scaling decisions

## Extensibility Points

### Middleware Integration

```go
// Add custom middleware to HTTP stack
func customMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom logic here
        next.ServeHTTP(w, r)
    })
}
```

### Logging Middleware

```go
// Add custom logging middleware
logger := logging.New("myapp",
    logging.WithProfile(logging.ProfileEnterprise),
    logging.WithMiddleware(
        logging.CorrelationMiddleware(),
        logging.RedactSecretsMiddleware(),
        customMiddleware(),
    ))
```

### Metrics Integration

```go
// Add custom metrics
customCounter := telemetry.Counter("custom_operations_total")
customGauge := telemetry.Gauge("custom_state")
customHistogram := telemetry.Histogram("custom_duration_seconds")
```

## Standards Compliance

### Fulmen Forge Workhorse Standard

- ✅ **App Identity**: `.fulmen/app.yaml` with dynamic configuration
- ✅ **Signal Handling**: Cross-platform graceful shutdown and reload
- ✅ **Exit Codes**: Standardized foundry exit codes
- ✅ **Observability**: Structured logging and metrics
- ✅ **Configuration**: Three-layer loading with XDG compliance
- ✅ **CLI Framework**: Standard commands and help system

### Crucible Integration

- ✅ **Schema Validation**: Embedded Crucible schemas for validation
- ✅ **Standards Access**: Programmatic access to Crucible standards
- ✅ **Version Alignment**: Synchronized with Crucible releases

## Resources

### Documentation

- [README.md](../../README.md) - Project overview and quick start
- [DEVELOPMENT.md](../../DEVELOPMENT.md) - Development handbook and workflows
- [CDRL Guide](fulmen_cdrl_guide.md) - Template customization guide
- [Accessing Crucible Docs](accessing-crucible-docs-via-gofulmen.md) - Embedded documentation access

### Standards & Specifications

- [Fulmen Forge Workhorse Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-workhorse-standard.md)
- [Agentic Attribution Standard](https://github.com/fulmenhq/crucible/blob/main/docs/standards/agentic-attribution.md)
- [Repository Safety Protocols](../../REPOSITORY_SAFETY_PROTOCOLS.md)

### External References

- [Chi Router Documentation](https://github.com/go-chi/chi)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Viper Configuration](https://github.com/spf13/viper)
- [Zap Logging](https://github.com/uber-go/zap)

## Version Information

- **Current Version**: 0.1.2
- **Gofulmen Version**: 0.1.15
- **Crucible Version**: 0.2.16
- **Go Version**: 1.25.1+
- **License**: MIT

## Contact & Support

- **Maintainer**: @3leapsdave (Dave Thompson)
- **AI Co-Maintainer**: ⚙️ Forge Architect (@forge-architect)
- **Issues**: [GitHub Issues](https://github.com/fulmenhq/forge-workhorse-groningen/issues)
- **Mattermost**: `#agents-groningen` (provisioning in progress)

---

_This document serves as the comprehensive overview of the Groningen template. For specific implementation details, see the package documentation and code comments._
