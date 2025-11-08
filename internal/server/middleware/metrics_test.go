package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/gofulmen/telemetry"
	"github.com/fulmenhq/gofulmen/telemetry/exporters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestMetrics_BasicFunctionality(t *testing.T) {
	// Setup real telemetry system for testing
	exporter := exporters.NewPrometheusExporter("test", ":0") // :0 for random port
	require.NoError(t, exporter.Start())
	defer exporter.Stop()

	config := &telemetry.Config{
		Enabled: true,
		Emitter: exporter,
	}

	sys, err := telemetry.NewSystem(config)
	require.NoError(t, err)

	// Replace global telemetry system temporarily
	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = sys
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with metrics middleware
	middleware := RequestMetrics(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Execute request
	middleware.ServeHTTP(rec, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test response", rec.Body.String())

	// Verify metrics are available in exporter
	metrics := exporter.GetMetrics()
	assert.NotEmpty(t, metrics, "Metrics should be recorded")
}

func TestRequestMetrics_WithTelemetryDisabled(t *testing.T) {
	// Disable telemetry
	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = nil
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestMetrics(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Should not panic and should work normally
	middleware.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequestMetrics_WithErrorStatus(t *testing.T) {
	// Setup real telemetry system
	exporter := exporters.NewPrometheusExporter("test", ":0")
	require.NoError(t, exporter.Start())
	defer exporter.Stop()

	config := &telemetry.Config{
		Enabled: true,
		Emitter: exporter,
	}

	sys, err := telemetry.NewSystem(config)
	require.NoError(t, err)

	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = sys
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	middleware := RequestMetrics(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Should record metrics including error counter
	metrics := exporter.GetMetrics()
	assert.NotEmpty(t, metrics, "Metrics should be recorded for error responses")
}

func TestRequestMetrics_WithRequestSize(t *testing.T) {
	// Setup real telemetry system
	exporter := exporters.NewPrometheusExporter("test", ":0")
	require.NoError(t, exporter.Start())
	defer exporter.Stop()

	config := &telemetry.Config{
		Enabled: true,
		Emitter: exporter,
	}

	sys, err := telemetry.NewSystem(config)
	require.NoError(t, err)

	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = sys
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestMetrics(handler)

	// Create request with content length
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Length", "1024")
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Should record request size metric
	metrics := exporter.GetMetrics()
	assert.NotEmpty(t, metrics, "Metrics should be recorded including request size")
}

func TestRequestMetrics_WithResponseSize(t *testing.T) {
	// Setup real telemetry system
	exporter := exporters.NewPrometheusExporter("test", ":0")
	require.NoError(t, exporter.Start())
	defer exporter.Stop()

	config := &telemetry.Config{
		Enabled: true,
		Emitter: exporter,
	}

	sys, err := telemetry.NewSystem(config)
	require.NoError(t, err)

	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = sys
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	responseBody := "test response with some content"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	})

	middleware := RequestMetrics(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Should record response size metric
	metrics := exporter.GetMetrics()
	assert.NotEmpty(t, metrics, "Metrics should be recorded including response size")
}

func TestGetEndpointPattern_StandardPaths(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/health", "/health/*"},
		{"/health/live", "/health/*"},
		{"/health/ready", "/health/*"},
		{"/health/startup", "/health/*"},
		{"/version", "/version"},
		{"/metrics", "/metrics"},
		{"/api/users/123", "/unknown"},
		{"/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			pattern := getEndpointPattern(req)
			assert.Equal(t, tt.expected, pattern, "Path %s should map to pattern %s", tt.path, tt.expected)
		})
	}
}

func TestRequestMetrics_WithRequestID(t *testing.T) {
	// Setup real telemetry system
	exporter := exporters.NewPrometheusExporter("test", ":0")
	require.NoError(t, exporter.Start())
	defer exporter.Stop()

	config := &telemetry.Config{
		Enabled: true,
		Emitter: exporter,
	}

	sys, err := telemetry.NewSystem(config)
	require.NoError(t, err)

	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = sys
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply RequestID middleware first, then RequestMetrics
	middleware := RequestID(RequestMetrics(handler))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Request ID should be in response header
	assert.Equal(t, "test-request-id", rec.Header().Get("X-Request-ID"))

	// Metrics should be recorded
	metrics := exporter.GetMetrics()
	assert.NotEmpty(t, metrics, "Metrics should be recorded")
}

func TestRequestMetrics_DurationMeasurement(t *testing.T) {
	// Setup real telemetry system
	exporter := exporters.NewPrometheusExporter("test", ":0")
	require.NoError(t, exporter.Start())
	defer exporter.Stop()

	config := &telemetry.Config{
		Enabled: true,
		Emitter: exporter,
	}

	sys, err := telemetry.NewSystem(config)
	require.NoError(t, err)

	originalTelemetry := observability.TelemetrySystem
	observability.TelemetrySystem = sys
	defer func() {
		observability.TelemetrySystem = originalTelemetry
	}()

	// Handler with artificial delay
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Small delay
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestMetrics(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	start := time.Now()
	middleware.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	// Should record duration metric
	metrics := exporter.GetMetrics()
	assert.NotEmpty(t, metrics, "Duration metrics should be recorded")
	assert.True(t, elapsed >= 10*time.Millisecond, "Should have waited at least 10ms")
}
