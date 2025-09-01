package testcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var toolsCmd = &cobra.Command{
	Use:     "list-tools",
	Aliases: []string{},
	Short:   "list tools",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		err := listToolsCmd(ctx, os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	testCmd.AddCommand(toolsCmd)
}

func listToolsCmd(ctx context.Context, cmd string) error {
	c, err := client.NewStdioMCPClient(cmd, nil, "mcp")
	if err != nil {
		return err
	}

	stdio := c.GetTransport().(*transport.Stdio)
	defer c.Close()

	_, err = c.Initialize(ctx, initReq)
	if err != nil {
		return err
	}
	go func() {
		copyStderrLogs(stdio.Stderr())
	}()

	result, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return err
	}

	log.Printf("result []tool:")
	for _, tool := range result.Tools {
		bs, _ := json.MarshalIndent(tool, "", "  ")
		log.Printf("tool %s:", tool.Name)
		fmt.Println(string(bs))
	}

	return nil
}
