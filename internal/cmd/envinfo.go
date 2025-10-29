package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fulmenhq/gofulmen/crucible"
)

var envInfoCmd = &cobra.Command{
	Use:   "envinfo",
	Short: "Display environment information",
	Long:  "Display comprehensive environment, configuration, and version information.",
	Run: func(cmd *cobra.Command, args []string) {
		version := crucible.GetVersion()

		fmt.Println("=== Groningen Environment Information ===\n")

		// Application Info
		fmt.Println("Application:")
		fmt.Printf("  Name:       groningen\n")
		fmt.Printf("  Version:    %s\n", versionInfo.Version)
		fmt.Printf("  Commit:     %s\n", versionInfo.Commit)
		fmt.Printf("  Built:      %s\n", versionInfo.BuildDate)
		fmt.Println()

		// SSOT Info
		fmt.Println("SSOT:")
		fmt.Printf("  Gofulmen:   %s\n", version.Gofulmen)
		fmt.Printf("  Crucible:   %s\n", version.Crucible)
		fmt.Println()

		// Runtime Info
		fmt.Println("Runtime:")
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  GOOS:       %s\n", runtime.GOOS)
		fmt.Printf("  GOARCH:     %s\n", runtime.GOARCH)
		fmt.Printf("  NumCPU:     %d\n", runtime.NumCPU())
		fmt.Println()

		// Configuration
		fmt.Println("Configuration:")
		fmt.Printf("  Server Host:    %s\n", viper.GetString("server.host"))
		fmt.Printf("  Server Port:    %d\n", viper.GetInt("server.port"))
		fmt.Printf("  Log Level:      %s\n", viper.GetString("logging.level"))
		fmt.Printf("  Log Profile:    %s\n", viper.GetString("logging.profile"))
		fmt.Printf("  Metrics Port:   %d\n", viper.GetInt("metrics.port"))
		fmt.Printf("  Config File:    %s\n", viper.ConfigFileUsed())
		if viper.ConfigFileUsed() == "" {
			fmt.Printf("                  (using defaults and environment variables)\n")
		}
		fmt.Println()

		fmt.Println("=== End Environment Information ===")
	},
}

func init() {
	rootCmd.AddCommand(envInfoCmd)
}
