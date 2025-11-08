package observability

import (
	"fmt"

	"github.com/fulmenhq/gofulmen/telemetry"
	"github.com/fulmenhq/gofulmen/telemetry/exporters"
)

var (
	// TelemetrySystem is the global telemetry system
	TelemetrySystem *telemetry.System

	// PrometheusExporter is the prometheus metrics exporter
	PrometheusExporter *exporters.PrometheusExporter
)

// InitMetrics initializes the telemetry system with Prometheus exporter.
// The exporter listens on the provided port (use 0 for random assignment).
// Optional namespace parameter for telemetry integration.
func InitMetrics(serviceName string, port int, namespace ...string) error {
	if port <= 0 {
		port = 9090
	}

	// Use namespace if provided, otherwise use service name
	metricNamespace := serviceName
	if len(namespace) > 0 && namespace[0] != "" {
		metricNamespace = namespace[0]
	}

	endpoint := fmt.Sprintf(":%d", port)

	// Create Prometheus exporter with namespace
	PrometheusExporter = exporters.NewPrometheusExporter(metricNamespace, endpoint)

	// Start Prometheus HTTP server
	if err := PrometheusExporter.Start(); err != nil {
		return err
	}

	// Create telemetry system with Prometheus exporter
	config := &telemetry.Config{
		Enabled: true,
		Emitter: PrometheusExporter,
	}

	sys, err := telemetry.NewSystem(config)
	if err != nil {
		return err
	}

	TelemetrySystem = sys
	return nil
}
