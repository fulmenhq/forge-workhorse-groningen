package observability

import (
	"github.com/fulmenhq/gofulmen/telemetry"
	"github.com/fulmenhq/gofulmen/telemetry/exporters"
)

var (
	// TelemetrySystem is the global telemetry system
	TelemetrySystem *telemetry.System

	// PrometheusExporter is the prometheus metrics exporter
	PrometheusExporter *exporters.PrometheusExporter
)

// InitMetrics initializes the telemetry system with Prometheus exporter
// Optional namespace parameter for telemetry integration
func InitMetrics(serviceName string, namespace ...string) error {
	// Use namespace if provided, otherwise use service name
	metricNamespace := serviceName
	if len(namespace) > 0 && namespace[0] != "" {
		metricNamespace = namespace[0]
	}

	// Create Prometheus exporter with namespace
	PrometheusExporter = exporters.NewPrometheusExporter(metricNamespace, ":9090")

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
