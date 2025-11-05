package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fulmenhq/gofulmen/pkg/signals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
)

var (
	serverPort int
	serverHost string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long: `Start the HTTP server with graceful shutdown support.

Signal Handling:
  • Ctrl+C (SIGINT) or SIGTERM: Graceful shutdown
  • Ctrl+C twice within 2s: Force quit
  • SIGHUP: Config reload (placeholder - restart recommended)

The server will cleanly shut down the HTTP server and flush logs on shutdown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get app identity for telemetry namespace
		identity := GetAppIdentity()
		namespace := identity.TelemetryNamespace()

		// Initialize server logger with namespace
		logLevel := viper.GetString("logging.level")
		observability.InitServerLogger(identity.BinaryName, logLevel, namespace)

		// Initialize metrics with namespace
		if err := observability.InitMetrics(identity.BinaryName, namespace); err != nil {
			observability.ServerLogger.Error("Failed to initialize metrics",
				zap.Error(err))
			return fmt.Errorf("metrics initialization failed: %w", err)
		}

		observability.ServerLogger.Info("Initializing server",
			zap.String("service", identity.BinaryName),
			zap.String("namespace", namespace),
			zap.String("version", versionInfo.Version),
			zap.String("host", serverHost),
			zap.Int("port", serverPort))

		// Create server
		srv := server.New(serverHost, serverPort)

		// Set app identity for handlers
		handlers.SetAppIdentity(identity)

		// Get shutdown timeout from config
		shutdownTimeout := viper.GetDuration("server.shutdown_timeout")
		if shutdownTimeout == 0 {
			shutdownTimeout = 10 * time.Second
		}

		// Register graceful shutdown handlers (LIFO order - last registered, first executed)
		// Handler 1: Flush logger (executed last)
		signals.OnShutdown(func(ctx context.Context) error {
			observability.ServerLogger.Info("Flushing logger...")
			if err := observability.ServerLogger.Sync(); err != nil {
				// Sync errors are often benign (stdout/stderr already closed)
				observability.ServerLogger.Warn("Logger sync returned error (may be benign)",
					zap.Error(err))
			}
			return nil
		})

		// Handler 2: Shutdown HTTP server (executed first)
		signals.OnShutdown(func(ctx context.Context) error {
			observability.ServerLogger.Info("Shutting down HTTP server...")
			shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}

			observability.ServerLogger.Info("HTTP server stopped gracefully")
			return nil
		})

		// Register config reload handler (SIGHUP)
		signals.OnReload(func(ctx context.Context) error {
			observability.ServerLogger.Info("Received SIGHUP: config reload requested")
			// Note: Config reload implementation deferred - would need thread-safe config update
			observability.ServerLogger.Warn("Config reload not yet implemented - restart server to apply changes")
			return nil
		})

		// Enable double-tap force quit (Ctrl+C within 2 seconds)
		if err := signals.EnableDoubleTap(signals.DoubleTapConfig{
			Window:  2 * time.Second,
			Message: "Press Ctrl+C again within 2 seconds to force quit",
		}); err != nil {
			observability.ServerLogger.Warn("Failed to enable double-tap force quit",
				zap.Error(err))
		}

		// Start server in background goroutine
		errChan := make(chan error, 1)
		go func() {
			observability.ServerLogger.Info("Starting HTTP server...",
				zap.String("host", serverHost),
				zap.Int("port", serverPort))
			if err := srv.Start(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}()

		// Start signal listener in background
		go func() {
			if err := signals.Listen(cmd.Context()); err != nil {
				observability.ServerLogger.Error("Signal handler error", zap.Error(err))
				errChan <- err
			}
		}()

		// Wait for error or shutdown completion
		if err := <-errChan; err != nil {
			return fmt.Errorf("server error: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&serverHost, "host", "localhost", "server host")
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "server port")

	_ = viper.BindPFlag("server.host", serveCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
}
