package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents the aggregate health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// ProbeResponse represents individual probe response
type ProbeResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthChecker defines interface for health checkable components
type HealthChecker interface {
	CheckHealth(ctx context.Context) error
}

// HealthManager manages health checks and probe states
type HealthManager struct {
	checkers map[string]HealthChecker
	version  string
}

// NewHealthManager creates a new health manager
func NewHealthManager(version string) *HealthManager {
	return &HealthManager{
		checkers: make(map[string]HealthChecker),
		version:  version,
	}
}

// RegisterChecker registers a health checker
func (hm *HealthManager) RegisterChecker(name string, checker HealthChecker) {
	hm.checkers[name] = checker
}

// runHealthChecks executes all registered health checks
func (hm *HealthManager) runHealthChecks(ctx context.Context) map[string]string {
	checks := make(map[string]string)

	for name, checker := range hm.checkers {
		select {
		case <-ctx.Done():
			checks[name] = "timeout"
			return checks
		default:
			if err := checker.CheckHealth(ctx); err != nil {
				checks[name] = "unhealthy"
			} else {
				checks[name] = "healthy"
			}
		}
	}

	return checks
}

// determineOverallStatus determines overall health status
func (hm *HealthManager) determineOverallStatus(checks map[string]string) string {
	degraded := false
	for _, status := range checks {
		if status == "unhealthy" {
			return "unhealthy"
		}
		if status == "degraded" || status == "timeout" {
			degraded = true
		}
	}

	// If we recorded any degraded/timeout checks, reflect that in aggregate status
	if degraded {
		return "degraded"
	}

	return "healthy"
}

// HealthHandler handles aggregate health check requests
func (hm *HealthManager) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Run health checks with timeout
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	checks := hm.runHealthChecks(checkCtx)
	status := hm.determineOverallStatus(checks)

	response := HealthResponse{
		Status:    status,
		Version:   hm.version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	// Return appropriate status code
	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// LivenessHandler handles liveness probe requests
// Liveness indicates if the application is running
func (hm *HealthManager) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Run health checks with timeout for liveness (shorter timeout)
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	checks := hm.runHealthChecks(checkCtx)
	status := hm.determineOverallStatus(checks)

	response := ProbeResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
	}

	// Return appropriate status code
	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// ReadinessHandler handles readiness probe requests
// Readiness indicates if the application is ready to serve traffic
func (hm *HealthManager) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Run health checks with timeout for readiness
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	checks := hm.runHealthChecks(checkCtx)
	status := hm.determineOverallStatus(checks)

	response := ProbeResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
	}

	// Return appropriate status code
	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// StartupHandler handles startup probe requests
// Startup indicates if the application has completed initialization
func (hm *HealthManager) StartupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Run health checks with timeout for startup
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	checks := hm.runHealthChecks(checkCtx)
	status := hm.determineOverallStatus(checks)

	response := ProbeResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
	}

	// Return appropriate status code
	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// Global health manager instance
var globalHealthManager *HealthManager

// InitHealthManager initializes the global health manager
func InitHealthManager(version string) {
	globalHealthManager = NewHealthManager(version)
}

// GetHealthManager returns the global health manager
func GetHealthManager() *HealthManager {
	return globalHealthManager
}

// LivenessHandler is the backward-compatible handler that uses the global manager
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	if globalHealthManager != nil {
		globalHealthManager.LivenessHandler(w, r)
	} else {
		// Fallback if not initialized
		response := ProbeResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}

// ReadinessHandler is the backward-compatible handler that uses the global manager
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if globalHealthManager != nil {
		globalHealthManager.ReadinessHandler(w, r)
	} else {
		// Fallback if not initialized
		response := ProbeResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}

// StartupHandler is the backward-compatible handler that uses the global manager
func StartupHandler(w http.ResponseWriter, r *http.Request) {
	if globalHealthManager != nil {
		globalHealthManager.StartupHandler(w, r)
	} else {
		// Fallback if not initialized
		response := ProbeResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}

// HealthHandler is the backward-compatible handler that uses the global manager
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if globalHealthManager != nil {
		globalHealthManager.HealthHandler(w, r)
	} else {
		// Fallback if not initialized
		response := HealthResponse{
			Status:    "healthy",
			Version:   AppVersion,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
