package tools

import (
	"context"
	"fmt"
	"github.com/mozillazg/kube-audit-mcp/pkg/utils"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mozillazg/kube-audit-mcp/pkg/config"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
)

type QueryAuditLogTool struct {
	params types.QueryAuditLogParams
	cfg    *config.Config
}

var resourceMapping = map[string]string{
	"pod":                "pods",
	"deployment":         "deployments",
	"service":            "services",
	"svc":                "services",
	"configmap":          "configmaps",
	"cm":                 "configmaps",
	"secret":             "secrets",
	"sec":                "secrets",
	"role":               "roles",
	"rolebinding":        "rolebindings",
	"clusterrole":        "clusterroles",
	"clusterrolebinding": "clusterrolebindings",
	"node":               "nodes",
	"namespace":          "namespaces",
	"ns":                 "namespaces",
	"pv":                 "persistentvolumes",
	"pvc":                "persistentvolumeclaims",
	"sa":                 "serviceaccounts",
	"deploy":             "deployments",
	"rs":                 "replicasets",
	"ds":                 "daemonsets",
	"sts":                "statefulsets",
	"ing":                "ingresses",
}

const auditLogResultNote = `Notes:
- When the 'ImpersonatedUser' field in an audit log is not empty, 
  it indicates the operation was performed via Kubernetes' user impersonation feature.
  In this scenario, the 'user' field represents the impersonated identity under which the action was executed, 
  while the 'ImpersonatedUser' field identifies the actual user who initiated the request.
  Therefore, to audit the true actor, you must refer to the 'ImpersonatedUser' field, if it is present in the log entry.
`

func NewQueryAuditLogTool(cfg *config.Config) *QueryAuditLogTool {
	return &QueryAuditLogTool{cfg: cfg}
}

func (t *QueryAuditLogTool) Register(s *server.MCPServer) {
	s.AddTool(t.newTool(), t.handle)
}

func (t *QueryAuditLogTool) handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var input types.QueryAuditLogParams
	if err := req.BindArguments(&input); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	input = t.normalizeParams(input)
	p, err := t.cfg.GetProviderByName(input.ClusterName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := p.QueryAuditLog(ctx, input)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result.Params = input
	result.Note = auditLogResultNote

	return mcp.NewToolResultStructuredOnly(result), nil
}

func (t *QueryAuditLogTool) normalizeParams(params types.QueryAuditLogParams) types.QueryAuditLogParams {
	if params.ClusterName == "" {
		params.ClusterName = t.cfg.DefaultCluster
	}
	if params.StartTime.IsZero() {
		params.StartTime = types.NewTimeParam(time.Now().UTC().Add(-24 * time.Hour * 7))
	}
	if params.EndTime.IsZero() {
		params.EndTime = types.NewTimeParam(time.Now().UTC())
	}
	if params.Limit <= 0 {
		params.Limit = 10
	} else if params.Limit > 20 {
		params.Limit = 20
	}
	newResourceTypes := make([]string, 0, len(params.ResourceTypes))
	for _, rt := range params.ResourceTypes {
		if rt == "" {
			continue
		}
		rt = strings.ToLower(rt)
		if mapped, ok := resourceMapping[rt]; ok {
			newResourceTypes = append(newResourceTypes, mapped)
		} else {
			newResourceTypes = append(newResourceTypes, rt)
		}
	}
	params.ResourceTypes = newResourceTypes
	if len(params.ResourceTypes) > 0 {
		params.ResourceTypes = utils.RemoveDuplicates(params.ResourceTypes)
	}

	if len(params.Verbs) > 0 {
		if utils.Contains(params.Verbs, "update") {
			params.Verbs = append(params.Verbs, "create", "patch", "delete", "deletecollection")
		}
		params.Verbs = utils.RemoveDuplicates(params.Verbs)
	}

	return params
}

func (t *QueryAuditLogTool) newTool() mcp.Tool {
	return mcp.NewTool("query_audit_log",
		mcp.WithDescription(`Query Kubernetes (k8s) audit logs.`),
		mcp.WithString("namespace",
			mcp.Description(`(Optional) Match by namespace. 

Supports exact matching and suffix wildcards:
- Exact match: "default", "kube-system", "kube-public"
- Suffix wildcard: "kube*", "app-*" (matches namespaces that start with the specified prefix)
`),
		),
		mcp.WithArray("verbs",
			mcp.Description(`(Optional) Filter by action verbs, multiple values are allowed.

Common values:
- "get": Get a resource
- "list": List resources
- "create": Create a resource
- "update": Update a resource
- "delete": Delete a resource
- "patch": Partially update a resource
- "watch": Watch for changes to a resource
`),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithArray("resource_types",
			mcp.Description(`(Optional) K8s resource type, multiple values are allowed.

Supports full names and short names. Common values:
- Core Resources: pods(pod), services(svc), configmaps(cm), secrets, nodes, namespaces(ns)
- Application Resources: deployments(deploy), replicasets(rs), daemonsets(ds), statefulsets(sts)
- Storage Resources: persistentvolumes(pv), persistentvolumeclaims(pvc)
- Network Resources: ingresses(ing), networkpolicies
- RBAC Resources: roles, rolebindings, clusterroles, clusterrolebindings

If you are uncertain about the resource type, you can call the 'list_common_resource_types()' tool to view common resource types
or ask the user to provide the corresponding one.
`),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithString("resource_name",
			mcp.Description(`(Optional) Match by resource name. 

Supports exact matching and suffix wildcards:
- Exact match: "nginx-deployment", "my-service"
- Suffix wildcard: "nginx-*", "app-*" (matches resource names that start with the specified prefix)
`),
		),
		mcp.WithString("user",
			mcp.Description(`(Optional) Match by user name. 

Supports exact matching and suffix wildcards:
- Exact match: "system:admin", "kubernetes-admin"
- Suffix wildcard: "system:*", "kube*" (matches users that start with the specified prefix)
`),
		),
		mcp.WithString("start_time",
			mcp.Description(`(Optional) Query start time. 

Supported formats:
- ISO 8601 format: "2024-01-01T10:00:00"
- Relative time: "30m" (30 minutes ago), "1h" (1 hour ago), "24h" (24 hours ago), "7d" (7 days ago)
- Defaults to "7d" (i.e., queries logs from the last 7 days).
`),
			mcp.DefaultString("7d"),
		),
		mcp.WithString("end_time",
			mcp.Description(`(Optional) Query end time. 

Supported formats:
- ISO 8601 format: "2024-01-01T10:00:00"
- Relative time: "30m" (30 minutes ago), "1h" (1 hour ago), "24h" (24 hours ago), "7d" (7 days ago)
- If empty, it defaults to the current time.
`),
		),
		mcp.WithNumber("limit",
			mcp.Description(`(Optional) Result limit, defaults to 10. Maximum is 20.`),
			mcp.Min(1),
			mcp.Max(20),
			mcp.DefaultNumber(10),
		),
		mcp.WithString("cluster_name",
			mcp.Description(fmt.Sprintf(`(Optional) The name of the cluster to query audit logs from.

You can use the 'list_clusters()' tool to view available clusters and their names,
If not specified, it defaults to the configured default cluster (%s).`, t.cfg.DefaultCluster)),
			mcp.DefaultString(t.cfg.DefaultCluster),
			mcp.Enum(t.cfg.AvailableClusterNames()...),
		),
	)
}
