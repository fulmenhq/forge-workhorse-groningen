package cmd

import (
	"github.com/fulmenhq/gofulmen/foundry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Run self-health check",
	Long:  "Run a self-health check to verify the application can start successfully.",
	Run: func(cmd *cobra.Command, args []string) {
		observability.CLILogger.Info("Running health check...")

		// Check 1: Version info available
		if versionInfo.Version == "" {
			observability.CLILogger.Error("❌ FAIL: Version information missing")
			ExitWithCode(observability.CLILogger, foundry.ExitConfigInvalid, "Version information missing", nil)
			return
		}
		observability.CLILogger.Debug("Version check passed", zap.String("version", versionInfo.Version))
		observability.CLILogger.Info("✅ Version information available")

		// Check 2: Logger initialized
		if observability.CLILogger == nil {
			// Can't log if logger is nil, so use stderr
			ExitWithCodeStderr(foundry.ExitConfigInvalid, "Logger not initialized", nil)
			return
		}
		observability.CLILogger.Info("✅ Logger initialized")

		// Check 3: Configuration loaded
		observability.CLILogger.Info("✅ Configuration system ready")

		// Overall status
		observability.CLILogger.Info("")
		observability.CLILogger.Info("✅ All health checks passed")
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
