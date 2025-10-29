package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/fulmenhq/gofulmen/crucible"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run diagnostic checks",
	Long:  "Run diagnostic checks on the system and suggest fixes for common issues.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Groningen Doctor ===\n")
		fmt.Println("Running diagnostic checks...\n")

		allChecks := true

		// Check 1: Go version
		fmt.Printf("[1/5] Checking Go version... ")
		goVersion := runtime.Version()
		if goVersion >= "go1.23" {
			fmt.Printf("✅ %s\n", goVersion)
		} else {
			fmt.Printf("⚠️  %s (recommended: go1.23+)\n", goVersion)
			allChecks = false
		}

		// Check 2: Crucible access
		fmt.Printf("[2/5] Checking Crucible access... ")
		version := crucible.GetVersion()
		if version.Crucible != "" {
			fmt.Printf("✅ v%s\n", version.Crucible)
		} else {
			fmt.Printf("❌ Cannot access Crucible\n")
			allChecks = false
		}

		// Check 3: Gofulmen access
		fmt.Printf("[3/5] Checking Gofulmen access... ")
		if version.Gofulmen != "" {
			fmt.Printf("✅ v%s\n", version.Gofulmen)
		} else {
			fmt.Printf("❌ Cannot access Gofulmen\n")
			allChecks = false
		}

		// Check 4: Config directory
		fmt.Printf("[4/5] Checking config directory... ")
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Printf("⚠️  Cannot find config directory\n")
			allChecks = false
		} else {
			fmt.Printf("✅ %s\n", configDir)
		}

		// Check 5: Environment
		fmt.Printf("[5/5] Checking environment... ")
		fmt.Printf("✅ %s/%s\n", runtime.GOOS, runtime.GOARCH)

		fmt.Println()
		if allChecks {
			fmt.Println("✅ All checks passed! Your groningen installation is healthy.")
		} else {
			fmt.Println("⚠️  Some checks failed. Review the output above for details.")
		}
		fmt.Println("\n=== End Diagnostics ===")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
