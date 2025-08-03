package tools

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListCommonResourceTypesTool struct{}

var commonResourceTypes = map[string][]string{
	"Core Resources": {
		"pods", "services", "endpoints", "persistentvolumes", "persistentvolumeclaims",
		"configmaps", "secrets", "nodes", "namespaces", "serviceaccounts",
	},
	"Apps Resources": {
		"deployments", "replicasets", "daemonsets", "statefulsets",
	},
	"Networking Resources": {
		"networkpolicies", "ingresses", "ingressclasses",
	},
	"Storage Resources": {
		"storageclasses", "volumeattachments", "csidrivers", "csinodes",
	},
	"RBAC Resources": {
		"roles", "rolebindings", "clusterroles", "clusterrolebindings",
	},
	"Extension Resources": {
		"customresourcedefinitions", "validatingadmissionwebhooks",
		"mutatingadmissionwebhooks",
	},
	"Scheduling Resources": {
		"priorityclasses", "poddisruptionbudgets",
	},
	"Monitoring Resources": {
		"events",
	},
}

func (t *ListCommonResourceTypesTool) Register(s *server.MCPServer) {
	s.AddTool(t.newTool(), t.handle)
}

func (t *ListCommonResourceTypesTool) handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resourceTypes := make(map[string][]string)
	for category, types := range commonResourceTypes {
		resourceTypes[category] = types
	}

	return mcp.NewToolResultStructuredOnly(resourceTypes), nil
}

func (t *ListCommonResourceTypesTool) newTool() mcp.Tool {
	return mcp.NewTool("list_common_resource_types",
		mcp.WithDescription(`列出常见的 Kubernetes 资源类型，帮助选择正确的 resource_type 参数。
    
返回按类别分组的常见K8s资源类型列表。`),
	)
}
