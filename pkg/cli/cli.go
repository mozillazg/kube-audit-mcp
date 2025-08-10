package cli

import (
	"github.com/spf13/cobra"
)

var onlyVersion bool

var (
	rootCmd = &cobra.Command{
		Use:   "kube-audit-mcp",
		Short: "MCP Server for Kubernetes Audit Logs.",
		Long: `MCP Server for Kubernetes Audit Logs.

More info: https://github.com/mozillazg/kube-audit-mcp`,
		Run: func(cmd *cobra.Command, args []string) {
			if onlyVersion {
				printVersion()
				return
			}
			_ = cmd.Help()
		},
	}
)

func init() {
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(sampleConfCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().BoolVar(&onlyVersion, "version", false, "show version")
}

func Run() error {
	return rootCmd.Execute()
}
