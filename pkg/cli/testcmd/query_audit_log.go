package testcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

type queryAuditLogOptions struct {
	types.QueryAuditLogParams

	startTime string
	endTime   string
	config    string
}

var queryAuditLogOpts = &queryAuditLogOptions{}

var queryAuditLogCmd = &cobra.Command{
	Use:     "query-audit-log",
	Aliases: []string{"query_audit_log", "query-auditlog", "query_auditlog", "queryauditlog"},
	Short:   "run query_audit_log",
	Run: func(cmd *cobra.Command, args []string) {
		err := runQueryAuditLogCmd(os.Args[0], queryAuditLogOpts)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	queryAuditLogCmd.Flags().StringVar(&queryAuditLogOpts.config, "config", "",
		"Path to the configuration file. If not specified, ~/.config/kube-audit-mcp/config.yaml will be used.")
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.ClusterName, "cluster-name",
		"",
		"Cluster name to query audit log.",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.startTime, "start-time",
		"7d",
		"Start time to query audit log. e.g. '10h', '1d'",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.endTime, "end-time",
		"1s",
		"End time to query audit log. e.g. '10h', '1d'",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.User, "user",
		"",
		"User to query audit log.",
	)
	queryAuditLogCmd.Flags().StringSliceVar(
		&queryAuditLogOpts.Verbs, "verb",
		nil,
		"Verbs to query audit log.",
	)
	queryAuditLogCmd.Flags().StringSliceVar(
		&queryAuditLogOpts.ResourceTypes, "resource-type",
		nil,
		"Resource types to query audit log. e.g. pods, deployments",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.ResourceName, "resource-name",
		"",
		"Resource name to query audit log.",
	)
	queryAuditLogCmd.Flags().IntVar(
		&queryAuditLogOpts.Limit, "limit",
		10,
		"Limit the number of results returned.",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.Namespace, "namespace",
		"",
		"Namespace to query audit log.",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.ResourceName, "name",
		"",
		"Name to query audit log.",
	)

	testCmd.AddCommand(queryAuditLogCmd)
}

func runQueryAuditLogCmd(cmd string, opts *queryAuditLogOptions) error {
	ctx := context.Background()
	args := []string{"mcp"}
	if opts.config != "" {
		args = append(args, "--config", opts.config)
	}
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
		copyLogs(stdio.Stderr())
	}()

	result, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "query_audit_log",
			Arguments: map[string]interface{}{
				"cluster_name":   opts.ClusterName,
				"start_time":     opts.startTime,
				"end_time":       opts.endTime,
				"user":           opts.User,
				"verbs":          opts.Verbs,
				"resource_types": opts.ResourceTypes,
				"resource_name":  opts.ResourceName,
				"limit":          opts.Limit,
				"namespace":      opts.Namespace,
			},
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

func copyLogs(r io.Reader) {
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
