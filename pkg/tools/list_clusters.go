package tools

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mozillazg/kube-audit-mcp/pkg/config"
)

type ListClustersTool struct {
	cfg *config.Config

	clusters ClustersResult
}

type ClustersResult struct {
	DefaultCluster string        `json:"default_cluster"`
	Clusters       []ClusterInfo `json:"clusters"`
}

type ClusterInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Alias       []string `json:"alias,omitempty"`
	Disabled    bool     `json:"disabled"`
	Provider    string   `json:"provider"`
}

func NewListClustersTool(cfg *config.Config) *ListClustersTool {
	result := ClustersResult{
		DefaultCluster: cfg.DefaultCluster,
	}
	for _, c := range cfg.Clusters {
		info := ClusterInfo{
			Name:        c.Name,
			Description: c.Description,
			Alias:       c.Alias,
			Disabled:    c.Disabled,
			Provider:    c.Provider.Name,
		}
		result.Clusters = append(result.Clusters, info)
	}

	return &ListClustersTool{
		cfg:      cfg,
		clusters: result,
	}
}

func (t *ListClustersTool) Register(s *server.MCPServer) {
	s.AddTool(t.newTool(), t.handle)
}

func (t *ListClustersTool) handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := t.clusters

	return mcp.NewToolResultStructuredOnly(result), nil
}

func (t *ListClustersTool) newTool() mcp.Tool {
	return mcp.NewTool("list_clusters",
		mcp.WithDescription(
			`List all configured clusters in the MCP server.`),
	)
}
