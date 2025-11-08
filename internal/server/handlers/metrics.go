package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var metricsProxyClient = &http.Client{
	Timeout: 5 * time.Second,
}

// MetricsHandler proxies Prometheus metrics from the internal exporter so callers
// can scrape /metrics on the main HTTP server.
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	exporter := observability.PrometheusExporter
	if exporter == nil {
		WriteServiceUnavailable(w, "Metrics exporter not initialized")
		return
	}

	// Build metrics URL using the configured metrics port
	metricsPort := viper.GetInt("metrics.port")
	if metricsPort == 0 {
		metricsPort = 9090
	}
	metricsURL := fmt.Sprintf("http://127.0.0.1:%d/metrics", metricsPort)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, metricsURL, nil)
	if err != nil {
		WriteServiceUnavailable(w, "Unable to construct metrics request")
		return
	}

	// Preserve caller hint for content negotiation
	if accept := r.Header.Get("Accept"); accept != "" {
		req.Header.Set("Accept", accept)
	}

	resp, err := metricsProxyClient.Do(req)
	if err != nil {
		if observability.ServerLogger != nil {
			observability.ServerLogger.Warn("Failed to proxy metrics request",
				zap.String("url", metricsURL),
				zap.Error(err))
		}
		WriteServiceUnavailable(w, "Prometheus exporter unavailable")
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			observability.ServerLogger.Warn("Failed to close metrics response body",
				zap.Error(err))
		}
	}()

	for key, values := range resp.Header {
		// Skip hop-by-hop headers; net/http handles them.
		if strings.EqualFold(key, "Connection") || strings.EqualFold(key, "Keep-Alive") ||
			strings.EqualFold(key, "Proxy-Authenticate") || strings.EqualFold(key, "Proxy-Authorization") ||
			strings.EqualFold(key, "TE") || strings.EqualFold(key, "Trailer") ||
			strings.EqualFold(key, "Transfer-Encoding") || strings.EqualFold(key, "Upgrade") {
			continue
		}

		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	// Ensure we always advertise Prometheus content type
	if resp.Header.Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil && observability.ServerLogger != nil {
		observability.ServerLogger.Warn("Failed to write metrics response",
			zap.Error(err))
	}
}
