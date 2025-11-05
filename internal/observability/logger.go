package observability

import (
	"fmt"
	"os"

	"github.com/fulmenhq/gofulmen/logging"
)

var (
	// CLILogger is used for CLI commands (SIMPLE profile)
	CLILogger *logging.Logger

	// ServerLogger is used for HTTP server (STRUCTURED profile)
	ServerLogger *logging.Logger
)

// InitCLILogger initializes the CLI logger with SIMPLE profile
func InitCLILogger(serviceName string, verbose bool) {
	// Use the simplified NewCLI helper for CLI logging
	logger, err := logging.NewCLI(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize CLI logger: %v\n", err)
		os.Exit(1)
	}

	// Set level to DEBUG if verbose
	if verbose {
		logger.SetLevel(logging.DEBUG)
	}

	CLILogger = logger
}

// InitServerLogger initializes the server logger with STRUCTURED profile
// Optional namespace parameter for telemetry integration
func InitServerLogger(serviceName string, logLevel string, namespace ...string) {
	level := parseLogLevel(logLevel)

	// Build static fields with optional namespace
	staticFields := make(map[string]any)
	if len(namespace) > 0 && namespace[0] != "" {
		staticFields["namespace"] = namespace[0]
	}

	config := &logging.LoggerConfig{
		Profile:      logging.ProfileStructured,
		DefaultLevel: level,
		Service:      serviceName,
		Environment:  "production",
		StaticFields: staticFields,
		Middleware: []logging.MiddlewareConfig{
			{
				Name:    "correlation",
				Enabled: true,
				Order:   100,
				Config:  make(map[string]any),
			},
		},
		Sinks: []logging.SinkConfig{
			{
				Type:   "console",
				Format: "json",
				Console: &logging.ConsoleSinkConfig{
					Stream:   "stderr",
					Colorize: false,
				},
			},
		},
		EnableCaller:     true,
		EnableStacktrace: true,
	}

	logger, err := logging.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize server logger: %v\n", err)
		os.Exit(1)
	}

	ServerLogger = logger
}

// parseLogLevel converts string log level to logging severity string
func parseLogLevel(levelStr string) string {
	switch levelStr {
	case "trace":
		return "TRACE"
	case "debug":
		return "DEBUG"
	case "info":
		return "INFO"
	case "warn", "warning":
		return "WARN"
	case "error":
		return "ERROR"
	default:
		return "INFO"
	}
}
