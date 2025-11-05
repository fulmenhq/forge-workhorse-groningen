package cmd

import (
	"fmt"
	"os"

	"github.com/fulmenhq/gofulmen/foundry"
	"github.com/fulmenhq/gofulmen/logging"
	"go.uber.org/zap"
)

// ExitWithCode exits the program with a semantic foundry exit code and logs the error.
// This helper ensures consistent error logging with exit code metadata before exiting.
//
// Parameters:
//   - logger: The logger to use for error output (can be nil for early failures)
//   - exitCode: The foundry exit code constant (e.g., foundry.ExitConfigInvalid)
//   - msg: Human-readable error message
//   - err: The underlying error (can be nil)
func ExitWithCode(logger *logging.Logger, exitCode foundry.ExitCode, msg string, err error) {
	// Get exit code metadata from foundry catalog
	info, ok := foundry.GetExitCodeInfo(exitCode)
	if !ok {
		// Fallback if we can't get exit code info (should never happen)
		fmt.Fprintf(os.Stderr, "FATAL: %s: %v (exit code: %d)\n", msg, err, exitCode)
		os.Exit(int(exitCode))
	}

	// Log error with exit code metadata
	if logger != nil {
		// Use structured logger if available
		logger.Error(msg,
			zap.Error(err),
			zap.Int("exit_code", info.Code),
			zap.String("exit_name", info.Name),
			zap.String("exit_description", info.Description),
			zap.String("exit_category", info.Category),
		)
	} else {
		// Fall back to stderr if no logger available
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %s: %v\n", msg, err)
		} else {
			fmt.Fprintf(os.Stderr, "FATAL: %s\n", msg)
		}
		fmt.Fprintf(os.Stderr, "Exit Code: %d (%s) - %s\n", info.Code, info.Name, info.Description)
	}

	// Exit with semantic code
	os.Exit(info.Code)
}

// ExitWithCodeStderr is a variant that writes to stderr without a logger.
// Use this for early failures before logger initialization.
//
// Parameters:
//   - exitCode: The foundry exit code constant
//   - msg: Human-readable error message
//   - err: The underlying error (can be nil)
func ExitWithCodeStderr(exitCode foundry.ExitCode, msg string, err error) {
	info, ok := foundry.GetExitCodeInfo(exitCode)
	if !ok {
		// Fallback if we can't get exit code info
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %s: %v (exit code: %d)\n", msg, err, exitCode)
		} else {
			fmt.Fprintf(os.Stderr, "FATAL: %s (exit code: %d)\n", msg, exitCode)
		}
		os.Exit(int(exitCode))
	}

	// Write to stderr with exit code metadata
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s: %v\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", msg)
	}
	fmt.Fprintf(os.Stderr, "Exit Code: %d (%s) - %s\n", info.Code, info.Name, info.Description)

	os.Exit(info.Code)
}
