package testcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"log"
)

func callTool(ctx context.Context, cmd string, args []string,
	toolName string, arguments map[string]interface{}) error {
	c, err := client.NewStdioMCPClient(cmd, nil, args...)
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

	result, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: arguments,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("result []content:")
	for _, content := range result.Content {
		switch content.(type) {
		case mcp.TextContent:
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(content.(mcp.TextContent).Text), &data); err == nil {
				bs, _ := json.MarshalIndent(data, "", "  ")
				log.Print("TextContent:")
				fmt.Println(string(bs))
			} else {
				log.Print("TextContent:")
				fmt.Println(content.(mcp.TextContent).Text)
			}
			break
		default:
			log.Printf("Unknown content type: %T", content)
		}
	}

	return nil
}

func copyStderrLogs(r io.Reader) {
	buf := make([]byte, 1024*128)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			log.Printf("stderr: %s", string(buf[:n]))
		}
		if err != nil {
			break
		}
	}
}
