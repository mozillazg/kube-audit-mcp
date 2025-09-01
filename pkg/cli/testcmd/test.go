package testcmd

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{},
	Short:   "test tools",
}

var initReq = mcp.InitializeRequest{
	Params: mcp.InitializeParams{
		ProtocolVersion: "2024-11-05",
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.Implementation{
			Name:    "My Application",
			Version: "1.0.0",
		},
	},
}

func Registry(root *cobra.Command) {
	root.AddCommand(testCmd)
}
