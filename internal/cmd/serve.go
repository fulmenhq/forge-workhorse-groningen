package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server"
)

var (
	serverPort int
	serverHost string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  "Start the HTTP server with graceful shutdown support. Use Ctrl+C to trigger graceful shutdown.",
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

		// Start server in goroutine
		errChan := make(chan error, 1)
		go func() {
			if err := srv.Start(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}()

		// Wait for interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		select {
		case err := <-errChan:
			return fmt.Errorf("server error: %w", err)
		case sig := <-sigChan:
			observability.ServerLogger.Info("Received shutdown signal",
				zap.String("signal", sig.String()))

			// Graceful shutdown with timeout
			shutdownTimeout := viper.GetDuration("server.shutdown_timeout")
			if shutdownTimeout == 0 {
				shutdownTimeout = 10 * time.Second
			}

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}

			observability.ServerLogger.Info("Server stopped gracefully")
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
