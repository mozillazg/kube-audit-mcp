package cli

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "kube-audit-mcp",
		Short: "MCP Server for Kubernetes Audit Logs.",
		Long: `MCP Server for Kubernetes Audit Logs.

More info: https://github.com/mozillazg/kube-audit-mcp`,
	}
)

func init() {
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(sampleConfCmd)
}

func Run() error {
	return rootCmd.Execute()
}
