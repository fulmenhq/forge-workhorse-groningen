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
func InitMetrics(serviceName string) error {
	// Create Prometheus exporter
	PrometheusExporter = exporters.NewPrometheusExporter(serviceName, ":9090")

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
