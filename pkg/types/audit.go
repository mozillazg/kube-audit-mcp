package types

import k8saudit "k8s.io/apiserver/pkg/apis/audit"

type AuditLogEntry k8saudit.Event

type AuditLogResult struct {
	Entries       []AuditLogEntry     `json:"entries"`
	Total         int                 `json:"total"`
	ProviderQuery string              `json:"provider_query"`
	Params        QueryAuditLogParams `json:"params"`
}
