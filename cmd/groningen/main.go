package main

import (
	"fmt"
	"os"

	"github.com/fulmenhq/gofulmen/foundry"

	"github.com/fulmenhq/forge-workhorse-groningen/internal/cmd"
	"github.com/fulmenhq/forge-workhorse-groningen/internal/server/handlers"
)

// Version information set via ldflags during build
// Example: go build -ldflags="-X main.version=1.0.0 -X main.commit=abc123 -X main.buildDate=2025-10-28"
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	// Set version info for commands to access
	cmd.SetVersionInfo(version, commit, buildDate)

	// Set version info for HTTP handlers
	handlers.SetVersionInfo(version, commit, buildDate)

	// Execute root command
	if err := cmd.Execute(); err != nil {
		// Command execution failed - use generic failure code
		// Individual commands may have already logged specific errors
		info, _ := foundry.GetExitCodeInfo(foundry.ExitFailure)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Exit Code: %d (%s)\n", info.Code, info.Name)
		os.Exit(info.Code)
	}
}
