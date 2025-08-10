package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

// This will be overridden at build time
var (
	version = "unknown"
	commit  = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of kube-audit-mcp.",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("Version:   %s\n", version)
	fmt.Printf("GitCommit: %s\n", commit)
	fmt.Println("For more information, visit https://github.com/mozillazg/kube-audit-mcp")
}
