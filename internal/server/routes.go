package server

import (
	"os"

	"github.com/fulmenhq/gofulmen/pkg/signals"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
)

// registerRoutes registers all HTTP routes
func (s *Server) registerRoutes() {
	// Health endpoint
	s.router.Get("/health", handlers.HealthHandler)

	// Version endpoint
	s.router.Get("/version", handlers.VersionHandler)

	// Metrics endpoint
	s.router.Get("/metrics", handlers.MetricsHandler)

	// Admin signal endpoint (optional, requires GRONINGEN_ADMIN_TOKEN)
	s.registerAdminEndpoint()
}

// registerAdminEndpoint optionally registers the admin signal endpoint
func (s *Server) registerAdminEndpoint() {
	// Get admin token from environment (identity-aware)
	adminToken := os.Getenv("GRONINGEN_ADMIN_TOKEN")
	if adminToken == "" {
		observability.ServerLogger.Debug("Admin signal endpoint disabled (no GRONINGEN_ADMIN_TOKEN set)")
		return
	}

	// Create HTTP signal handler with bearer token auth and rate limiting
	handler := signals.NewHTTPHandler(signals.HTTPConfig{
		TokenAuth: adminToken,
		RateLimit: 10,  // 10 requests per minute
		RateBurst: 5,   // burst size
		Manager:   nil, // use default global manager
	})

	// Register admin endpoint
	s.router.Post("/admin/signal", handler.ServeHTTP)

	observability.ServerLogger.Info("Admin signal endpoint enabled",
		zap.String("path", "/admin/signal"),
		zap.String("auth", "bearer token"),
		zap.String("rate_limit", "10/min, burst 5"))
	observability.ServerLogger.Warn("Admin endpoint enabled - ensure this server is not exposed to public internet")
}
