package provider

import (
	"context"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
)

type Provider interface {
	QueryAuditLog(context.Context, types.QueryAuditLogParams) (types.AuditLogResult, error)
}
