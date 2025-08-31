package gcp

import (
	"testing"

	"github.com/mozillazg/kube-audit-mcp/pkg/types"
)

func TestCloudLoggingProvider_logName(t *testing.T) {
	type fields struct {
		projectId string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				projectId: "test",
			},
			want: "projects/test/logs/cloudaudit.googleapis.com%2Factivity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudLoggingProvider{
				projectId: tt.fields.projectId,
			}
			if got := c.logName(); got != tt.want {
				t.Errorf("logName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudLoggingProvider_buildQuery(t *testing.T) {
	type fields struct {
		projectId   string
		clusterName string
	}
	tests := []struct {
		name   string
		fields fields
		params types.QueryAuditLogParams
		want   string
	}{
		{
			name: "basic query with project id only",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{},
			want:   `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity"`,
		},
		{
			name: "query with cluster name",
			fields: fields{
				projectId:   "test-project",
				clusterName: "test-cluster",
			},
			params: types.QueryAuditLogParams{},
			want:   `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND resource.labels.cluster_name="test-cluster"`,
		},
		{
			name: "query with user filter",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				User: "test-user@example.com",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND protoPayload.authenticationInfo.principalEmail: "test-user@example.com"`,
		},
		{
			name: "query with namespace filter",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				Namespace: "test-namespace",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND protoPayload.resourceName: "/namespaces/test-namespace/"`,
		},
		{
			name: "query with verbs",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				Verbs: []string{"create", "delete"},
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND protoPayload.methodName: (".create" OR ".delete")`,
		},
		{
			name: "query with resource types",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				ResourceTypes: []string{"pods", "services"},
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND protoPayload.resourceName: ("/pods/" OR "/services/")`,
		},
		{
			name: "query with resource name",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				ResourceName: "test-pod",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND protoPayload.resourceName =~ "/test-pod$"`,
		},
		{
			name: "complex query with multiple filters",
			fields: fields{
				projectId:   "test-project",
				clusterName: "test-cluster",
			},
			params: types.QueryAuditLogParams{
				User:          "test-user@example.com",
				Namespace:     "test-namespace",
				Verbs:         []string{"create"},
				ResourceTypes: []string{"pods"},
				ResourceName:  "test-pod",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND resource.labels.cluster_name="test-cluster" AND protoPayload.authenticationInfo.principalEmail: "test-user@example.com" AND protoPayload.resourceName: "/namespaces/test-namespace/" AND protoPayload.methodName: (".create") AND protoPayload.resourceName: ("/pods/") AND protoPayload.resourceName =~ "/test-pod$"`,
		},
		{
			name: "query with wildcard user should not add user filter",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				User: "*",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity"`,
		},
		{
			name: "query with wildcard namespace should not add namespace filter",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				Namespace: "*",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity"`,
		},
		{
			name: "query with wildcard resource name should not add resource name filter",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				ResourceName: "*",
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity"`,
		},
		{
			name: "query with empty arrays should not add filters",
			fields: fields{
				projectId: "test-project",
			},
			params: types.QueryAuditLogParams{
				Verbs:         []string{},
				ResourceTypes: []string{},
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudLoggingProvider{
				projectId:   tt.fields.projectId,
				clusterName: tt.fields.clusterName,
			}
			if got := c.buildQuery(tt.params); got != tt.want {
				t.Errorf("buildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
