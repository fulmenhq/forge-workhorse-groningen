package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
)

// MetricsResponse represents basic metrics information
type MetricsResponse struct {
	Status           string `json:"status"`
	MetricsAvailable bool   `json:"metricsAvailable"`
	PrometheusPort   int    `json:"prometheusPort"`
	Message          string `json:"message"`
}

// MetricsHandler serves basic metrics information
// Full Prometheus metrics are available on port 9090/metrics
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	response := MetricsResponse{
		Status:           "ok",
		MetricsAvailable: observability.TelemetrySystem != nil,
		PrometheusPort:   9090,
		Message:          "Full Prometheus metrics available at http://localhost:9090/metrics",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
