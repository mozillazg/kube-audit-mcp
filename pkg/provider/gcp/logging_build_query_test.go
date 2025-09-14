package gcp

import (
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloudLoggingProvider_buildQuery2(t *testing.T) {
	type fields struct {
		client      cloudLoggingProviderClientInterface
		projectId   string
		clusterName string
	}
	type args struct {
		params types.QueryAuditLogParams
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "should build a query with a namespace and wildcard",
			fields: fields{
				projectId:   "test-project",
				clusterName: "test-cluster",
			},
			args: args{
				params: types.QueryAuditLogParams{
					Namespace: "test-namespace*",
				},
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND resource.labels.cluster_name="test-cluster" AND protoPayload.resourceName: "/namespaces/test-namespace"`,
		},
		{
			name: "should build a query with a resource name and wildcard",
			fields: fields{
				projectId:   "test-project",
				clusterName: "test-cluster",
			},
			args: args{
				params: types.QueryAuditLogParams{
					ResourceName: "test-resource*",
				},
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND resource.labels.cluster_name="test-cluster" AND protoPayload.resourceName =~ "/test-resource"`,
		},
		{
			name: "should build a query with a user name and wildcard",
			fields: fields{
				projectId:   "test-project",
				clusterName: "test-cluster",
			},
			args: args{
				params: types.QueryAuditLogParams{
					User: "test-user*",
				},
			},
			want: `resource.type="k8s_cluster" AND logName="projects/test-project/logs/cloudaudit.googleapis.com%2Factivity" AND resource.labels.cluster_name="test-cluster" AND protoPayload.authenticationInfo.principalEmail: "test-user"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudLoggingProvider{
				client:      tt.fields.client,
				projectId:   tt.fields.projectId,
				clusterName: tt.fields.clusterName,
			}
			assert.Equalf(t, tt.want, c.buildQuery(tt.args.params), "buildQuery(%v)", tt.args.params)
		})
	}
}
