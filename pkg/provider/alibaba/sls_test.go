package alibaba

import (
	"testing"
	"time"

	"github.com/mozillazg/kube-audit-mcp/pkg/types"
)

func TestSLSProvider_buildQuery(t *testing.T) {
	provider := &SLSProvider{}

	tests := []struct {
		name     string
		params   types.QueryAuditLogParams
		expected string
	}{
		{
			name: "empty params - should return base query",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				Limit:     100,
			},
			expected: "*",
		},
		{
			name: "user parameter",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				User:      "testuser",
				Limit:     100,
			},
			expected: `* and user.username: "testuser"`,
		},
		{
			name: "user wildcard - should be ignored",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				User:      "*",
				Limit:     100,
			},
			expected: "*",
		},
		{
			name: "namespace parameter",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				Namespace: "kube-system",
				Limit:     100,
			},
			expected: `* and objectRef.namespace: "kube-system"`,
		},
		{
			name: "namespace wildcard - should be ignored",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				Namespace: "*",
				Limit:     100,
			},
			expected: "*",
		},
		{
			name: "single verb",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				Verbs:     []string{"get"},
				Limit:     100,
			},
			expected: `* and (verb: "get")`,
		},
		{
			name: "multiple verbs",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				Verbs:     []string{"get", "list", "create"},
				Limit:     100,
			},
			expected: `* and (verb: "get" or verb: "list" or verb: "create")`,
		},
		{
			name: "single resource type",
			params: types.QueryAuditLogParams{
				StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:       types.NewTimeParam(time.Now()),
				ResourceTypes: []string{"pods"},
				Limit:         100,
			},
			expected: `* and (objectRef.resource: "pods")`,
		},
		{
			name: "multiple resource types",
			params: types.QueryAuditLogParams{
				StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:       types.NewTimeParam(time.Now()),
				ResourceTypes: []string{"pods", "services", "deployments"},
				Limit:         100,
			},
			expected: `* and (objectRef.resource: "pods" or objectRef.resource: "services" or objectRef.resource: "deployments")`,
		},
		{
			name: "resource name",
			params: types.QueryAuditLogParams{
				StartTime:    types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:      types.NewTimeParam(time.Now()),
				ResourceName: "my-pod",
				Limit:        100,
			},
			expected: `* and objectRef.name: "my-pod"`,
		},
		{
			name: "resource name wildcard - should be ignored",
			params: types.QueryAuditLogParams{
				StartTime:    types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:      types.NewTimeParam(time.Now()),
				ResourceName: "*",
				Limit:        100,
			},
			expected: "*",
		},
		{
			name: "all parameters combined",
			params: types.QueryAuditLogParams{
				StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:       types.NewTimeParam(time.Now()),
				User:          "admin",
				Namespace:     "default",
				Verbs:         []string{"create", "update"},
				ResourceTypes: []string{"pods", "services"},
				ResourceName:  "my-resource",
				Limit:         100,
			},
			expected: `* and user.username: "admin" and objectRef.namespace: "default" and (verb: "create" or verb: "update") and (objectRef.resource: "pods" or objectRef.resource: "services") and objectRef.name: "my-resource"`,
		},
		{
			name: "mixed wildcards and values",
			params: types.QueryAuditLogParams{
				StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:       types.NewTimeParam(time.Now()),
				User:          "*",
				Namespace:     "kube-system",
				Verbs:         []string{"get"},
				ResourceTypes: []string{"secrets"},
				ResourceName:  "*",
				Limit:         100,
			},
			expected: `* and objectRef.namespace: "kube-system" and (verb: "get") and (objectRef.resource: "secrets")`,
		},
		{
			name: "empty string parameters - should be ignored",
			params: types.QueryAuditLogParams{
				StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:       types.NewTimeParam(time.Now()),
				User:          "",
				Namespace:     "",
				Verbs:         []string{},
				ResourceTypes: []string{},
				ResourceName:  "",
				Limit:         100,
			},
			expected: "*",
		},
		{
			name: "parameters with special characters",
			params: types.QueryAuditLogParams{
				StartTime:    types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:      types.NewTimeParam(time.Now()),
				User:         "user@domain.com",
				Namespace:    "test-namespace",
				ResourceName: "my-pod-with-dashes",
				Limit:        100,
			},
			expected: `* and user.username: "user@domain.com" and objectRef.namespace: "test-namespace" and objectRef.name: "my-pod-with-dashes"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.buildQuery(tt.params)
			if result != tt.expected {
				t.Errorf("buildQuery() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSLSProvider_buildQuery_EdgeCases(t *testing.T) {
	provider := &SLSProvider{}

	t.Run("nil verbs slice", func(t *testing.T) {
		params := types.QueryAuditLogParams{
			StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
			EndTime:   types.NewTimeParam(time.Now()),
			Verbs:     nil,
			Limit:     100,
		}
		result := provider.buildQuery(params)
		expected := "*"
		if result != expected {
			t.Errorf("buildQuery() with nil verbs = %q, want %q", result, expected)
		}
	})

	t.Run("nil resource types slice", func(t *testing.T) {
		params := types.QueryAuditLogParams{
			StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
			EndTime:       types.NewTimeParam(time.Now()),
			ResourceTypes: nil,
			Limit:         100,
		}
		result := provider.buildQuery(params)
		expected := "*"
		if result != expected {
			t.Errorf("buildQuery() with nil resource types = %q, want %q", result, expected)
		}
	})

	t.Run("empty verbs in slice", func(t *testing.T) {
		params := types.QueryAuditLogParams{
			StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
			EndTime:   types.NewTimeParam(time.Now()),
			Verbs:     []string{"", "get", ""},
			Limit:     100,
		}
		result := provider.buildQuery(params)
		expected := `* and (verb: "" or verb: "get" or verb: "")`
		if result != expected {
			t.Errorf("buildQuery() with empty verbs in slice = %q, want %q", result, expected)
		}
	})

	t.Run("empty resource types in slice", func(t *testing.T) {
		params := types.QueryAuditLogParams{
			StartTime:     types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
			EndTime:       types.NewTimeParam(time.Now()),
			ResourceTypes: []string{"", "pods", ""},
			Limit:         100,
		}
		result := provider.buildQuery(params)
		expected := `* and (objectRef.resource: "" or objectRef.resource: "pods" or objectRef.resource: "")`
		if result != expected {
			t.Errorf("buildQuery() with empty resource types in slice = %q, want %q", result, expected)
		}
	})
}
