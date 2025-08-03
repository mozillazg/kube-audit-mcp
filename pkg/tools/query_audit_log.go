package tools

import (
	"context"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
)

type QueryAuditLogTool struct {
	params types.QueryAuditLogParams
	p      provider.Provider
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

func NewQueryAuditLogTool(p provider.Provider) *QueryAuditLogTool {
	return &QueryAuditLogTool{p: p}
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

	logs, err := t.p.QueryAuditLog(ctx, input)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultStructuredOnly(logs), nil
}

func (t *QueryAuditLogTool) normalizeParams(params types.QueryAuditLogParams) types.QueryAuditLogParams {
	if params.StartTime.IsZero() {
		params.StartTime = types.NewTimeParam(time.Now().Add(-24 * time.Hour))
	}
	if params.EndTime.IsZero() {
		params.EndTime = types.NewTimeParam(time.Now())
	}
	if params.Limit <= 0 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
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
	return params
}

func (t *QueryAuditLogTool) newTool() mcp.Tool {
	return mcp.NewTool("query_audit_log",
		mcp.WithDescription(`查询 Kubernetes (k8s) 审计日志。

功能说明：
- 支持多种时间格式（ISO 8601 和相对时间）
- 支持资源类型的简写和模糊匹配
- 提供详细的参数验证和错误提示

使用建议：
1. 不确定资源类型时，可调用 list_common_resource_types() 查看常见资源类型或者询问用户提供对应的资源类型
2. 默认查询最近24小时的审计日志。默认限制返回10条记录
`),
		mcp.WithString("namespace",
			mcp.Description(`(可选) 按命名空间匹配。支持精确匹配和后缀通配符：
- 精确匹配: "default", "kube-system", "kube-public""
- 后缀通配符: "kube*", "app-*" (匹配以指定前缀开头的命名空间)
`),
		),
		mcp.WithArray("verbs",
			mcp.Description(`(可选) 筛选操作动词，可选多个。常见值：
- "get": 获取资源
- "list": 列出资源
- "create": 创建资源
- "update": 更新资源
- "delete": 删除资源
- "patch": 部分更新资源
- "watch": 监听资源变化
`),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithArray("resource_types",
			mcp.Description(`(可选) K8s资源类型, 可选多个。支持完整名称和简写, 常见值：
- 核心资源: pods(pod), services(svc), configmaps(cm), secrets, nodes, namespaces(ns)
- 应用资源: deployments(deploy), replicasets(rs), daemonsets(ds), statefulsets(sts)
- 存储资源: persistentvolumes(pv), persistentvolumeclaims(pvc)
- 网络资源: ingresses(ing), networkpolicies
- RBAC资源: roles, rolebindings, clusterroles, clusterrolebindings
提示：可以使用 list_common_resource_types() 工具查看常见列表
`),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithString("resource_name",
			mcp.Description(`(可选) 按资源名称匹配。支持精确匹配和后缀通配符：
- 精确匹配: "nginx-deployment", "my-service"
- 后缀通配符: "nginx-*", "app-*" (匹配以指定前缀开头的资源名)
`),
		),
		mcp.WithString("user",
			mcp.Description(`(可选) 按用户名匹配。支持精确匹配和后缀通配符：
- 精确匹配: "system:admin", "kubernetes-admin"
- 后缀通配符: "system:*", "kube*" (匹配以指定前缀开头的用户)
`),
		),
		mcp.WithString("start_time",
			mcp.Description(`(可选) 查询开始时间。支持格式：
- ISO 8601格式: "2024-01-01T10:00:00"
- 相对时间: "30m"(30分钟前), "1h"(1小时前), "24h"(24小时前), "7d"(7天前)
- 默认值为 "24h"（即查询最近24小时的日志）
`),
			mcp.DefaultString("24h"),
		),
		mcp.WithString("end_time",
			mcp.Description(`(可选) 查询结束时间。支持格式：
- ISO 8601格式: "2024-01-01T10:00:00"
- 相对时间: "30m"(30分钟前), "1h"(1小时前), "24h"(24小时前), "7d"(7天前)
- 如果为空，默认使用当前时间
`),
		),
		mcp.WithNumber("limit",
			mcp.Description(`(可选) 返回结果限制，默认值为 10。最大值为 100。`),
			mcp.Min(1),
			mcp.Max(100),
			mcp.DefaultNumber(10),
		),
	)
}
