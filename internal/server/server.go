package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
	servermw "github.com/fulmenhq/forge-workhorse-groningen/internal/server/middleware"
)

// Server represents the HTTP server
type Server struct {
	router *chi.Mux
	server *http.Server
	host   string
	port   int
}

// New creates a new HTTP server instance
func New(host string, port int) *Server {
	r := chi.NewRouter()

	// Standard chi middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Our custom middleware (recovery, error handling, request ID, then metrics)
	r.Use(servermw.Recovery)
	r.Use(servermw.ErrorHandler)
	r.Use(servermw.RequestID)
	r.Use(servermw.RequestMetrics)

	// Standardized error responses
	r.NotFound(handlers.NotFoundHandler)
	r.MethodNotAllowed(handlers.MethodNotAllowedHandler)

	s := &Server{
		router: r,
		host:   host,
		port:   port,
	}

	// Register routes
	s.registerRoutes()

	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	observability.ServerLogger.Info("Starting HTTP server",
		zap.String("host", s.host),
		zap.Int("port", s.port),
		zap.String("addr", addr))

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	observability.ServerLogger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

// Handler exposes the underlying router for testing and instrumentation
func (s *Server) Handler() http.Handler {
	return s.router
}
