package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/observability"
)

var (
	cfgFile string
	verbose bool

	// Version info set by main package
	versionInfo struct {
		Version   string
		Commit    string
		BuildDate string
	}
)

// SetVersionInfo is called by main package to set version information
func SetVersionInfo(version, commit, buildDate string) {
	versionInfo.Version = version
	versionInfo.Commit = commit
	versionInfo.BuildDate = buildDate
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "groningen",
	Short: "A Fulmen workhorse application for robust, scalable backends",
	Long: `groningen is a production-ready workhorse application template from the
FulmenHQ ecosystem. Named after the Groningen horse breed known for strength
and toughness in heavy work.

This template provides:
- HTTP server with standard endpoints (/health, /version, /metrics)
- Structured logging with progressive profiles (via gofulmen)
- Three-layer configuration management
- Graceful shutdown and signal handling
- Observability and telemetry built-in
- CLI commands for server management and diagnostics

Use the subcommands to perform specific operations.`,
	SilenceUsage:  true,  // Don't show usage on errors
	SilenceErrors: true,  // We'll handle error output ourselves
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/groningen/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (sets log level to debug)")

	// Bind flags to viper
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Initialize CLI logger early so we can use it in config loading
	observability.InitCLILogger("groningen", verbose)

	if cfgFile != "" {
		// Use config file from flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find XDG config directory
		configDir, err := os.UserConfigDir()
		if err != nil {
			if verbose {
				observability.CLILogger.Warn("Could not find user config directory", zap.Error(err))
			}
			// Fall back to home directory
			home, err := os.UserHomeDir()
			if err != nil {
				observability.CLILogger.Error("Could not find home directory", zap.Error(err))
				os.Exit(1)
			}
			configDir = home
			viper.AddConfigPath(home)
			viper.SetConfigName(".groningen")
		} else {
			// Use XDG config directory
			appConfigDir := configDir + "/groningen"
			viper.AddConfigPath(appConfigDir)
			viper.SetConfigName("config")
		}

		// Also search in current directory
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
	}

	// Read in environment variables with prefix GRONINGEN_
	viper.SetEnvPrefix("GRONINGEN")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			observability.CLILogger.Debug("Using config file", zap.String("path", viper.ConfigFileUsed()))
		}
	} else {
		// It's OK if config file doesn't exist, we have defaults
		if verbose {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				observability.CLILogger.Debug("No config file found, using defaults and environment variables")
			} else {
				observability.CLILogger.Warn("Error reading config file", zap.Error(err))
			}
		}
	}

	// Set defaults
	setDefaults()
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.shutdown_timeout", "10s")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.profile", "structured")

	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.port", 9090)

	// Health check defaults
	viper.SetDefault("health.enabled", true)

	// Worker defaults
	viper.SetDefault("workers", 4)

	// Debug defaults
	viper.SetDefault("debug.enabled", false)
	viper.SetDefault("debug.pprof_enabled", false)
}
