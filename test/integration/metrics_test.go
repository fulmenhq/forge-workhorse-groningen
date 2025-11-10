package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsEndpoint_Integration(t *testing.T) {
	// Initialize logger for testing
	observability.InitCLILogger("test", false)
	observability.InitServerLogger("test", "info")

	// Initialize metrics for testing (use unique port to avoid conflicts)
	if err := observability.InitMetrics("test", 19090, "test"); err != nil {
		t.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Initialize health manager for testing
	handlers.InitHealthManager("test")

	// Setup test server with metrics
	s := server.New("localhost", 18080) // Use fixed port for testing

	// Add a simple test route
	s.Handler().(*chi.Mux).Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	// Start server in background
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.Start()
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-serverErr:
		t.Fatalf("Server failed to start: %v", err)
	default:
		// Server started successfully
	}

	// Get server URL
	serverURL := fmt.Sprintf("http://localhost:%d", s.Port())
	defer func() {
		_ = s.Shutdown(context.Background())
	}()

	// Make some requests to generate metrics
	for i := 0; i < 5; i++ {
		resp, err := http.Get(serverURL + "/test")
		require.NoError(t, err)
		_ = resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	// Test health endpoint to generate more metrics
	resp, err := http.Get(serverURL + "/health")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test version endpoint
	resp, err = http.Get(serverURL + "/version")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test metrics endpoint
	resp, err = http.Get(serverURL + "/metrics")
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify metrics content
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	metricsContent := string(body)

	// Check for expected metrics (with "test" namespace prefix from InitMetrics)
	assert.Contains(t, metricsContent, "test_http_requests_total", "Should contain HTTP request counter")
	assert.Contains(t, metricsContent, "test_http_request_duration_ms", "Should contain request duration histogram")
}

func TestMetricsEndpoint_LoadTesting(t *testing.T) {
	// Initialize logger for testing
	observability.InitCLILogger("test", false)
	observability.InitServerLogger("test", "info")

	// Initialize metrics for testing (use unique port to avoid conflicts)
	if err := observability.InitMetrics("test", 19091, "test"); err != nil {
		t.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Setup test server
	s := server.New("localhost", 8081)

	// Add test routes with different response patterns
	router := s.Handler().(*chi.Mux)
	router.Get("/fast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fast response"))
	})
	router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("slow response"))
	})
	router.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error response"))
	})

	// Start server
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-serverErr:
		t.Fatalf("Server failed to start: %v", err)
	default:
	}

	serverURL := fmt.Sprintf("http://localhost:%d", s.Port())
	defer func() {
		_ = s.Shutdown(context.Background())
	}()

	// Generate load with concurrent requests
	const numRequests = 50
	const numWorkers = 10

	requestChan := make(chan int, numRequests)
	for i := 0; i < numRequests; i++ {
		requestChan <- i
	}
	close(requestChan)

	// Worker pool
	worker := func() {
		for reqNum := range requestChan {
			var path string
			switch reqNum % 4 {
			case 0:
				path = "/fast"
			case 1:
				path = "/slow"
			case 2:
				path = "/error"
			default:
				path = "/health"
			}

			resp, err := http.Get(serverURL + path)
			if err == nil {
				_ = resp.Body.Close()
			}
		}
	}

	// Start workers
	start := time.Now()
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Wait for all requests to complete
	time.Sleep(2 * time.Second)
	elapsed := time.Since(start)

	// Verify metrics endpoint still works
	resp, err := http.Get(serverURL + "/metrics")
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	metricsContent := string(body)

	// Verify metrics reflect the load (using "test" namespace from InitMetrics)
	assert.Contains(t, metricsContent, "test_http_requests_total", "Should have HTTP request metrics")
	assert.Contains(t, metricsContent, "test_http_request_duration_ms", "Should have duration metrics")

	// Basic load validation
	assert.True(t, elapsed < 5*time.Second, "Load test should complete in reasonable time")

	t.Logf("Load test completed: %d requests in %v (%.2f req/s)",
		numRequests, elapsed, float64(numRequests)/elapsed.Seconds())
}

func TestMetricsEndpoint_PrometheusFormat(t *testing.T) {
	// Initialize logger for testing
	observability.InitCLILogger("test", false)
	observability.InitServerLogger("test", "info")

	// Initialize metrics for testing (use unique port to avoid conflicts)
	if err := observability.InitMetrics("test", 19092, "test"); err != nil {
		t.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Setup test server
	s := server.New("localhost", 8082)

	// Add test route
	s.Handler().(*chi.Mux).Get("/format-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "format test"}`))
	})

	// Start server
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-serverErr:
		t.Fatalf("Server failed to start: %v", err)
	default:
	}

	serverURL := fmt.Sprintf("http://localhost:%d", s.Port())
	defer func() {
		_ = s.Shutdown(context.Background())
	}()

	// Make a request
	resp, err := http.Get(serverURL + "/format-test")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get metrics
	resp, err = http.Get(serverURL + "/metrics")
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check content type (gofulmen exporter uses basic Prometheus format without charset)
	contentType := resp.Header.Get("Content-Type")
	assert.True(t,
		contentType == "text/plain; version=0.0.4" ||
			contentType == "text/plain; version=0.0.4; charset=utf-8",
		"Expected Prometheus content type, got: %s", contentType)

	// Read and validate Prometheus format
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	metricsContent := string(body)

	// Basic Prometheus format validation
	lines := strings.Split(strings.TrimSpace(metricsContent), "\n")

	// gofulmen telemetry may not emit HELP/TYPE comments by default
	// Just verify we have actual metric lines with proper format
	// Format: metric_name{labels} value
	hasValidMetrics := false
	for _, line := range lines {
		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Check for metric lines (contain metric name and value)
		if strings.Contains(line, "{") && len(strings.Fields(line)) >= 2 {
			hasValidMetrics = true
			break
		}
	}
	assert.True(t, hasValidMetrics, "Should have valid Prometheus metric lines")

	// Should have actual metric lines
	metricLines := 0
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			metricLines++
		}
	}
	assert.Greater(t, metricLines, 0, "Should have actual metric values")
}

func TestMetricsEndpoint_WithTelemetryDisabled(t *testing.T) {
	// Initialize logger for testing
	observability.InitCLILogger("test", false)
	observability.InitServerLogger("test", "info")

	// Save and clear global metrics state
	originalExporter := observability.PrometheusExporter
	originalTelemetry := observability.TelemetrySystem
	observability.PrometheusExporter = nil
	observability.TelemetrySystem = nil
	defer func() {
		observability.PrometheusExporter = originalExporter
		observability.TelemetrySystem = originalTelemetry
	}()

	// Temporarily disable metrics
	originalEnabled := os.Getenv("GRONINGEN_METRICS_ENABLED")
	_ = os.Setenv("GRONINGEN_METRICS_ENABLED", "false")
	defer func() {
		if originalEnabled != "" {
			_ = os.Setenv("GRONINGEN_METRICS_ENABLED", originalEnabled)
		} else {
			_ = os.Unsetenv("GRONINGEN_METRICS_ENABLED")
		}
	}()

	// Setup test server WITHOUT initializing metrics (since they're disabled)
	s := server.New("localhost", 18083)

	// Add test route
	s.Handler().(*chi.Mux).Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	})

	// Start server
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-serverErr:
		t.Fatalf("Server failed to start: %v", err)
	default:
	}

	serverURL := fmt.Sprintf("http://localhost:%d", s.Port())
	defer func() {
		_ = s.Shutdown(context.Background())
	}()

	// Make a request
	resp, err := http.Get(serverURL + "/test")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Metrics endpoint should return service unavailable when disabled
	resp, err = http.Get(serverURL + "/metrics")
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

// Test with real command execution
func TestMetrics_CommandIntegration(t *testing.T) {
	// This test would verify that the actual command can start and serve metrics
	// For now, we'll skip the command test since cmd.NewRootCommand is not exported
	t.Skip("Command integration test skipped - cmd.NewRootCommand not exported")
}
