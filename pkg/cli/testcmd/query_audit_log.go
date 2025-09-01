package testcmd

import (
	"context"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	"github.com/spf13/cobra"
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
	Short:   "call query_audit_log",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		err := runQueryAuditLogCmd(ctx, os.Args[0], queryAuditLogOpts)
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
		"",
		"Start time to query audit log. e.g. '10h', '1d'",
	)
	queryAuditLogCmd.Flags().StringVar(
		&queryAuditLogOpts.endTime, "end-time",
		"",
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
		0,
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

func runQueryAuditLogCmd(ctx context.Context, cmd string, opts *queryAuditLogOptions) error {
	args := []string{"mcp"}
	if opts.config != "" {
		args = append(args, "--config", opts.config)
	}

	err := callTool(ctx, cmd, args, "query_audit_log", map[string]interface{}{
		"cluster_name":   opts.ClusterName,
		"start_time":     opts.startTime,
		"end_time":       opts.endTime,
		"user":           opts.User,
		"verbs":          opts.Verbs,
		"resource_types": opts.ResourceTypes,
		"resource_name":  opts.ResourceName,
		"limit":          opts.Limit,
		"namespace":      opts.Namespace,
	})

	return err
}
