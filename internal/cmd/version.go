package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/fulmenhq/gofulmen/crucible"
)

var extended bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information. Use --extended for full details including Crucible and Go versions.",
	Run: func(cmd *cobra.Command, args []string) {
		if extended {
			// Extended version output
			fmt.Printf("groningen %s\n", versionInfo.Version)
			fmt.Printf("Commit: %s\n", versionInfo.Commit)
			fmt.Printf("Built: %s\n", versionInfo.BuildDate)
			fmt.Printf("Go: %s\n", runtime.Version())
			fmt.Printf("\n")

			// Gofulmen and Crucible versions
			version := crucible.GetVersion()
			fmt.Printf("Gofulmen: %s\n", version.Gofulmen)
			fmt.Printf("Crucible: %s\n", version.Crucible)
		} else {
			// Basic version output
			fmt.Printf("groningen %s\n", versionInfo.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&extended, "extended", "e", false, "show extended version information")
}
