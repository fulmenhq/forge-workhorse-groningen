package server

import (
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
)

// registerRoutes registers all HTTP routes
func (s *Server) registerRoutes() {
	// Health endpoint
	s.router.Get("/health", handlers.HealthHandler)

	// Version endpoint
	s.router.Get("/version", handlers.VersionHandler)
}
