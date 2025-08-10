package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mozillazg/kube-audit-mcp/pkg/config"
	"github.com/mozillazg/kube-audit-mcp/pkg/tools"
	"github.com/spf13/cobra"
)

type Options struct {
	config    string
	transport string
	addr      string
}

var opts Options

var (
	mcpCmd = &cobra.Command{
		Use:   "mcp",
		Short: "MCP Server for Kubernetes Audit Logs.",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := runMcpServer(opts)
			if err != nil && errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		},
	}
)

func init() {
	mcpCmd.Flags().StringVarP(
		&opts.config, "config", "c",
		config.ShortHomePath(config.DefaultConfigFile()),
		"Path to the configuration file.")

	mcpCmd.Flags().StringVarP(
		&opts.transport, "transport", "t",
		"stdio", "Transport type for MCP server (stdio).")
	//mcpCmd.Flags().StringVarP(
	//	&opts.addr, "address", "s",
	//	"127.0.0.1:8081", "Address to listen on for SSE Transport.")
}

func runMcpServer(opts Options) error {
	cfgPath, err := config.ExpandPath(opts.config)
	if err != nil {
		return fmt.Errorf("expanding config path: %+v", err)
	}
	cfg, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		return fmt.Errorf("loading configuration: %+v", err)
	}
	if err := cfg.Init(); err != nil {
		return fmt.Errorf("initializing configuration: %+v", err)
	}

	s := server.NewMCPServer("kube-audit", version,
		server.WithToolCapabilities(true),
	)

	queryAuditLog := tools.NewQueryAuditLogTool(cfg)
	queryAuditLog.Register(s)
	listCommonResourceTypes := tools.ListCommonResourceTypesTool{}
	listCommonResourceTypes.Register(s)
	listClusters := tools.NewListClustersTool(cfg)
	listClusters.Register(s)

	switch opts.transport {
	//case "sse":
	//	log.Printf("Starting MCP server with SSE transport on %s", opts.addr)
	//	sseServer := server.NewSSEServer(s)
	//	if err := sseServer.Start(opts.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		return err
	//	}
	//	log.Print("MCP server with SSE transport stopped")
	//	break
	default:
		return server.ServeStdio(s)
	}

	return nil
}
