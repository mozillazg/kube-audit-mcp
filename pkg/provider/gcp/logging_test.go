package gcp

import (
	"testing"

	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	"google.golang.org/protobuf/types/known/structpb"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8saudit "k8s.io/apiserver/pkg/apis/audit"
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

func TestCloudLoggingProvider_getLevelAndStage(t *testing.T) {
	// Helper function to create a structpb.Struct for testing
	createTestStruct := func() *structpb.Struct {
		s, _ := structpb.NewStruct(map[string]interface{}{
			"test": "data",
		})
		return s
	}

	tests := []struct {
		name       string
		objRef     *k8saudit.ObjectReference
		status     *v1.Status
		req        *structpb.Struct
		resp       *structpb.Struct
		wantLevel  k8saudit.Level
		wantStage  k8saudit.Stage
		wantStatus *v1.Status // Expected status after modification
	}{
		{
			name: "attach subresource - should set special handling",
			objRef: &k8saudit.ObjectReference{
				APIGroup:    "v1",
				APIVersion:  "v1",
				Resource:    "pods",
				Name:        "test-pod",
				Namespace:   "test-ns",
				Subresource: "attach",
			},
			status: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
			req:       createTestStruct(),
			resp:      createTestStruct(),
			wantLevel: k8saudit.LevelRequest,
			wantStage: k8saudit.StageResponseStarted,
			wantStatus: &v1.Status{
				Status:  "Switching Protocols (inferred)",
				Code:    101,
				Message: "Switching Protocols (inferred)",
			},
		},
		{
			name: "exec subresource - should set special handling",
			objRef: &k8saudit.ObjectReference{
				APIGroup:    "v1",
				APIVersion:  "v1",
				Resource:    "pods",
				Name:        "test-pod",
				Namespace:   "test-ns",
				Subresource: "exec",
			},
			status: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
			req:       createTestStruct(),
			resp:      createTestStruct(),
			wantLevel: k8saudit.LevelRequest,
			wantStage: k8saudit.StageResponseStarted,
			wantStatus: &v1.Status{
				Status:  "Switching Protocols (inferred)",
				Code:    101,
				Message: "Switching Protocols (inferred)",
			},
		},
		{
			name: "regular subresource - should not trigger special handling",
			objRef: &k8saudit.ObjectReference{
				APIGroup:    "v1",
				APIVersion:  "v1",
				Resource:    "pods",
				Name:        "test-pod",
				Namespace:   "test-ns",
				Subresource: "status",
			},
			status: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
			req:       createTestStruct(),
			resp:      createTestStruct(),
			wantLevel: k8saudit.LevelRequestResponse,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
		{
			name:      "nil objRef with both req and resp nil - should return LevelMetadata",
			objRef:    nil,
			status:    &v1.Status{Status: "OK", Code: 200, Message: "OK"},
			req:       nil,
			resp:      nil,
			wantLevel: k8saudit.LevelMetadata,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
		{
			name: "req not nil, resp nil - should return LevelRequest",
			objRef: &k8saudit.ObjectReference{
				APIGroup:   "v1",
				APIVersion: "v1",
				Resource:   "pods",
				Name:       "test-pod",
				Namespace:  "test-ns",
			},
			status:    &v1.Status{Status: "OK", Code: 200, Message: "OK"},
			req:       createTestStruct(),
			resp:      nil,
			wantLevel: k8saudit.LevelRequest,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
		{
			name: "req nil, resp not nil - should return LevelRequestResponse",
			objRef: &k8saudit.ObjectReference{
				APIGroup:   "v1",
				APIVersion: "v1",
				Resource:   "pods",
				Name:       "test-pod",
				Namespace:  "test-ns",
			},
			status:    &v1.Status{Status: "OK", Code: 200, Message: "OK"},
			req:       nil,
			resp:      createTestStruct(),
			wantLevel: k8saudit.LevelRequestResponse,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
		{
			name: "both req and resp not nil - should return LevelRequestResponse",
			objRef: &k8saudit.ObjectReference{
				APIGroup:   "v1",
				APIVersion: "v1",
				Resource:   "pods",
				Name:       "test-pod",
				Namespace:  "test-ns",
			},
			status:    &v1.Status{Status: "OK", Code: 200, Message: "OK"},
			req:       createTestStruct(),
			resp:      createTestStruct(),
			wantLevel: k8saudit.LevelRequestResponse,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
		{
			name:      "attach with nil objRef - should not trigger special handling",
			objRef:    nil,
			status:    &v1.Status{Status: "OK", Code: 200, Message: "OK"},
			req:       createTestStruct(),
			resp:      createTestStruct(),
			wantLevel: k8saudit.LevelRequestResponse,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
		{
			name: "objRef without subresource - should use normal logic",
			objRef: &k8saudit.ObjectReference{
				APIGroup:   "v1",
				APIVersion: "v1",
				Resource:   "pods",
				Name:       "test-pod",
				Namespace:  "test-ns",
			},
			status:    &v1.Status{Status: "OK", Code: 200, Message: "OK"},
			req:       nil,
			resp:      nil,
			wantLevel: k8saudit.LevelMetadata,
			wantStage: k8saudit.StageResponseComplete,
			wantStatus: &v1.Status{
				Status:  "OK",
				Code:    200,
				Message: "OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudLoggingProvider{}

			// Make a copy of the status to avoid modifying the original
			statusCopy := &v1.Status{
				Status:  tt.status.Status,
				Code:    tt.status.Code,
				Message: tt.status.Message,
			}

			gotLevel, gotStage := c.getLevelAndStage(tt.objRef, statusCopy, tt.req, tt.resp)

			if gotLevel != tt.wantLevel {
				t.Errorf("getLevelAndStage() level = %v, want %v", gotLevel, tt.wantLevel)
			}
			if gotStage != tt.wantStage {
				t.Errorf("getLevelAndStage() stage = %v, want %v", gotStage, tt.wantStage)
			}

			// Check if status was modified correctly
			if statusCopy.Code != tt.wantStatus.Code {
				t.Errorf("getLevelAndStage() modified status.Code = %v, want %v", statusCopy.Code, tt.wantStatus.Code)
			}
			if statusCopy.Status != tt.wantStatus.Status {
				t.Errorf("getLevelAndStage() modified status.Status = %v, want %v", statusCopy.Status, tt.wantStatus.Status)
			}
			if statusCopy.Message != tt.wantStatus.Message {
				t.Errorf("getLevelAndStage() modified status.Message = %v, want %v", statusCopy.Message, tt.wantStatus.Message)
			}
		})
	}
}

func TestCloudLoggingProvider_getObjectReference(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		want         *k8saudit.ObjectReference
	}{
		{
			name:         "empty resource name",
			resourceName: "",
			want:         nil,
		},
		{
			name:         "namespaced resource with name (6 parts)",
			resourceName: "apps/v1/namespaces/default/deployments/nginx-deployment",
			want: &k8saudit.ObjectReference{
				APIGroup:   "apps",
				APIVersion: "v1",
				Namespace:  "default",
				Resource:   "deployments",
				Name:       "nginx-deployment",
			},
		},
		{
			name:         "namespaced resource without name (5 parts)",
			resourceName: "apps/v1/namespaces/default/deployments",
			want: &k8saudit.ObjectReference{
				APIGroup:   "apps",
				APIVersion: "v1",
				Namespace:  "default",
				Resource:   "deployments",
			},
		},
		{
			name:         "cluster-scoped resource (4 parts)",
			resourceName: "rbac.authorization.k8s.io/v1/clusterroles/admin",
			want: &k8saudit.ObjectReference{
				APIGroup:   "rbac.authorization.k8s.io",
				APIVersion: "v1",
				Resource:   "clusterroles",
				Name:       "admin",
			},
		},
		{
			name:         "core api resource with namespace and name",
			resourceName: "/v1/namespaces/kube-system/pods/coredns-558bd4d5db-abcde",
			want: &k8saudit.ObjectReference{
				APIGroup:   "",
				APIVersion: "v1",
				Namespace:  "kube-system",
				Resource:   "pods",
				Name:       "coredns-558bd4d5db-abcde",
			},
		},
		{
			name:         "core api resource without namespace",
			resourceName: "/v1/namespaces/kube-system/pods",
			want: &k8saudit.ObjectReference{
				APIGroup:   "",
				APIVersion: "v1",
				Namespace:  "kube-system",
				Resource:   "pods",
			},
		},
		{
			name:         "namespaced subresource (7+ parts with namespace)",
			resourceName: "/v1/namespaces/default/pods/test-pod/exec",
			want: &k8saudit.ObjectReference{
				APIGroup:    "",
				APIVersion:  "v1",
				Namespace:   "default",
				Resource:    "pods",
				Name:        "test-pod",
				Subresource: "exec",
			},
		},
		{
			name:         "namespaced subresource with multiple subresource parts",
			resourceName: "/v1/namespaces/default/pods/test-pod/log/previous",
			want: &k8saudit.ObjectReference{
				APIGroup:    "",
				APIVersion:  "v1",
				Namespace:   "default",
				Resource:    "pods",
				Name:        "test-pod",
				Subresource: "log",
			},
		},
		{
			name:         "cluster-scoped subresource (5+ parts without namespace)",
			resourceName: "/v1/nodes/node-1/status",
			want: &k8saudit.ObjectReference{
				APIGroup:    "",
				APIVersion:  "v1",
				Resource:    "nodes",
				Name:        "node-1",
				Subresource: "status",
			},
		},
		{
			name:         "cluster-scoped subresource with multiple parts",
			resourceName: "/v1/nodes/node-1/proxy/stats",
			want: &k8saudit.ObjectReference{
				APIGroup:    "",
				APIVersion:  "v1",
				Resource:    "nodes",
				Name:        "node-1",
				Subresource: "proxy",
			},
		},
		{
			name:         "custom resource with namespace",
			resourceName: "custom.io/v1/namespaces/default/mycustomresources/my-cr",
			want: &k8saudit.ObjectReference{
				APIGroup:   "custom.io",
				APIVersion: "v1",
				Namespace:  "default",
				Resource:   "mycustomresources",
				Name:       "my-cr",
			},
		},
		{
			name:         "custom resource cluster-scoped",
			resourceName: "custom.io/v1/mycustomclusterresources/my-ccr",
			want: &k8saudit.ObjectReference{
				APIGroup:   "custom.io",
				APIVersion: "v1",
				Resource:   "mycustomclusterresources",
				Name:       "my-ccr",
			},
		},
		{
			name:         "attach subresource",
			resourceName: "/v1/namespaces/default/pods/test-pod/attach",
			want: &k8saudit.ObjectReference{
				APIGroup:    "",
				APIVersion:  "v1",
				Namespace:   "default",
				Resource:    "pods",
				Name:        "test-pod",
				Subresource: "attach",
			},
		},
		{
			name:         "metrics server resource",
			resourceName: "metrics.k8s.io/v1beta1/namespaces/default/pods",
			want: &k8saudit.ObjectReference{
				APIGroup:   "metrics.k8s.io",
				APIVersion: "v1beta1",
				Namespace:  "default",
				Resource:   "pods",
			},
		},
		{
			name:         "scale subresource",
			resourceName: "apps/v1/namespaces/default/deployments/nginx/scale",
			want: &k8saudit.ObjectReference{
				APIGroup:    "apps",
				APIVersion:  "v1",
				Namespace:   "default",
				Resource:    "deployments",
				Name:        "nginx",
				Subresource: "scale",
			},
		},
		{
			name:         "service account token subresource",
			resourceName: "/v1/namespaces/default/serviceaccounts/default/token",
			want: &k8saudit.ObjectReference{
				APIGroup:    "",
				APIVersion:  "v1",
				Namespace:   "default",
				Resource:    "serviceaccounts",
				Name:        "default",
				Subresource: "token",
			},
		},
		{
			name:         "unparseable resource name (3 parts)",
			resourceName: "invalid/resource/name",
			want:         nil,
		},
		{
			name:         "unparseable resource name (2 parts)",
			resourceName: "invalid/resource",
			want:         nil,
		},
		{
			name:         "single part resource name",
			resourceName: "invalid",
			want:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudLoggingProvider{}
			got := c.getObjectReference(tt.resourceName)

			if got == nil && tt.want == nil {
				return
			}

			if got == nil || tt.want == nil {
				t.Errorf("getObjectReference() = %v, want %v", got, tt.want)
				return
			}

			if got.APIGroup != tt.want.APIGroup {
				t.Errorf("getObjectReference().APIGroup = %v, want %v", got.APIGroup, tt.want.APIGroup)
			}
			if got.APIVersion != tt.want.APIVersion {
				t.Errorf("getObjectReference().APIVersion = %v, want %v", got.APIVersion, tt.want.APIVersion)
			}
			if got.Namespace != tt.want.Namespace {
				t.Errorf("getObjectReference().Namespace = %v, want %v", got.Namespace, tt.want.Namespace)
			}
			if got.Resource != tt.want.Resource {
				t.Errorf("getObjectReference().Resource = %v, want %v", got.Resource, tt.want.Resource)
			}
			if got.Name != tt.want.Name {
				t.Errorf("getObjectReference().Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Subresource != tt.want.Subresource {
				t.Errorf("getObjectReference().Subresource = %v, want %v", got.Subresource, tt.want.Subresource)
			}
		})
	}
}
