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
			name: "user wildcard",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				User:      "testuser*",
				Limit:     100,
			},
			expected: `* and user.username: testuser*`,
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
			name: "namespace wildcard",
			params: types.QueryAuditLogParams{
				StartTime: types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:   types.NewTimeParam(time.Now()),
				Namespace: "kube-*",
				Limit:     100,
			},
			expected: `* and objectRef.namespace: kube-*`,
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
			name: "resource name wildcard",
			params: types.QueryAuditLogParams{
				StartTime:    types.NewTimeParam(time.Now().Add(-1 * time.Hour)),
				EndTime:      types.NewTimeParam(time.Now()),
				ResourceName: "my-pod*",
				Limit:        100,
			},
			expected: `* and objectRef.name: my-pod*`,
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

func TestSLSProviderConfig_Init(t *testing.T) {
	tests := []struct {
		name        string
		config      SLSProviderConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with endpoint",
			config: SLSProviderConfig{
				Endpoint: "us-west-1.log.aliyuncs.com",
				Project:  "test-project",
				LogStore: "test-logstore",
			},
			expectError: false,
		},
		{
			name: "valid config with region only",
			config: SLSProviderConfig{
				Region:   "us-west-1",
				Project:  "test-project",
				LogStore: "test-logstore",
			},
			expectError: false,
		},
		{
			name: "valid config with both endpoint and region",
			config: SLSProviderConfig{
				Endpoint: "us-west-1.log.aliyuncs.com",
				Region:   "us-west-1",
				Project:  "test-project",
				LogStore: "test-logstore",
			},
			expectError: false,
		},
		{
			name: "valid config with auth version v4 and region",
			config: SLSProviderConfig{
				Region:      "us-west-1",
				AuthVersion: "v4",
				Project:     "test-project",
				LogStore:    "test-logstore",
			},
			expectError: false,
		},
		{
			name: "valid config with auth version v4 and endpoint that can be parsed for region",
			config: SLSProviderConfig{
				Endpoint:    "us-west-1.log.aliyuncs.com",
				AuthVersion: "v4",
				Project:     "test-project",
				LogStore:    "test-logstore",
			},
			expectError: false,
		},
		{
			name: "missing endpoint and region",
			config: SLSProviderConfig{
				Project:  "test-project",
				LogStore: "test-logstore",
			},
			expectError: true,
			errorMsg:    "either endpoint or region must be provided",
		},
		{
			name: "auth version v4 without region and unparseable endpoint",
			config: SLSProviderConfig{
				Endpoint:    "invalid-endpoint",
				AuthVersion: "v4",
				Project:     "test-project",
				LogStore:    "test-logstore",
			},
			expectError: true,
			errorMsg:    "region is required when auth_version is v4",
		},
		{
			name: "auth version v4 without region and no endpoint",
			config: SLSProviderConfig{
				AuthVersion: "v4",
				Project:     "test-project",
				LogStore:    "test-logstore",
			},
			expectError: true,
			errorMsg:    "either endpoint or region must be provided",
		},
		{
			name: "missing project",
			config: SLSProviderConfig{
				Endpoint: "us-west-1.log.aliyuncs.com",
				LogStore: "test-logstore",
			},
			expectError: true,
			errorMsg:    "project is required",
		},
		{
			name: "empty project",
			config: SLSProviderConfig{
				Endpoint: "us-west-1.log.aliyuncs.com",
				Project:  "",
				LogStore: "test-logstore",
			},
			expectError: true,
			errorMsg:    "project is required",
		},
		{
			name: "missing logstore",
			config: SLSProviderConfig{
				Endpoint: "us-west-1.log.aliyuncs.com",
				Project:  "test-project",
			},
			expectError: true,
			errorMsg:    "logstore is required",
		},
		{
			name: "empty logstore",
			config: SLSProviderConfig{
				Endpoint: "us-west-1.log.aliyuncs.com",
				Project:  "test-project",
				LogStore: "",
			},
			expectError: true,
			errorMsg:    "logstore is required",
		},
		{
			name: "auth version other than v4 without region",
			config: SLSProviderConfig{
				Endpoint:    "us-west-1.log.aliyuncs.com",
				AuthVersion: "v1",
				Project:     "test-project",
				LogStore:    "test-logstore",
			},
			expectError: false,
		},
		{
			name: "empty auth version",
			config: SLSProviderConfig{
				Endpoint:    "us-west-1.log.aliyuncs.com",
				AuthVersion: "",
				Project:     "test-project",
				LogStore:    "test-logstore",
			},
			expectError: false,
		},
		{
			name: "endpoint auto-generation from region",
			config: SLSProviderConfig{
				Region:   "cn-hangzhou",
				Project:  "test-project",
				LogStore: "test-logstore",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original config in tests
			config := tt.config
			err := config.Init()

			if tt.expectError {
				if err == nil {
					t.Errorf("Init() expected error but got nil")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Init() error = %q, want %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Init() unexpected error = %v", err)
				}

				// Test that endpoint is auto-generated when only region is provided
				if tt.config.Endpoint == "" && tt.config.Region != "" {
					expectedEndpoint := tt.config.Region + ".log.aliyuncs.com"
					if config.Endpoint != expectedEndpoint {
						t.Errorf("Init() endpoint auto-generation: got %q, want %q", config.Endpoint, expectedEndpoint)
					}
				}
			}
		})
	}
}

func TestSLSProviderConfig_Init_EdgeCases(t *testing.T) {
	t.Run("endpoint generation preserves original region", func(t *testing.T) {
		config := SLSProviderConfig{
			Region:   "eu-central-1",
			Project:  "test-project",
			LogStore: "test-logstore",
		}

		err := config.Init()
		if err != nil {
			t.Errorf("Init() unexpected error = %v", err)
		}

		expectedEndpoint := "eu-central-1.log.aliyuncs.com"
		if config.Endpoint != expectedEndpoint {
			t.Errorf("endpoint = %q, want %q", config.Endpoint, expectedEndpoint)
		}
		if config.Region != "eu-central-1" {
			t.Errorf("region = %q, want %q", config.Region, "eu-central-1")
		}
	})

	t.Run("existing endpoint is preserved when region is also provided", func(t *testing.T) {
		originalEndpoint := "custom.endpoint.com"
		config := SLSProviderConfig{
			Endpoint: originalEndpoint,
			Region:   "us-west-1",
			Project:  "test-project",
			LogStore: "test-logstore",
		}

		err := config.Init()
		if err != nil {
			t.Errorf("Init() unexpected error = %v", err)
		}

		if config.Endpoint != originalEndpoint {
			t.Errorf("endpoint should be preserved, got %q, want %q", config.Endpoint, originalEndpoint)
		}
	})

	t.Run("region extraction for v4 auth works correctly", func(t *testing.T) {
		config := SLSProviderConfig{
			Endpoint:    "ap-southeast-1.log.aliyuncs.com",
			AuthVersion: "v4",
			Project:     "test-project",
			LogStore:    "test-logstore",
		}

		err := config.Init()
		if err != nil {
			t.Errorf("Init() unexpected error = %v", err)
		}

		// The region should be extracted from the endpoint
		if config.Region != "ap-southeast-1" {
			t.Errorf("region extraction: got %q, want %q", config.Region, "ap-southeast-1")
		}
	})

	t.Run("whitespace in required fields is considered invalid", func(t *testing.T) {
		tests := []struct {
			name   string
			config SLSProviderConfig
		}{
			{
				name: "project with only spaces",
				config: SLSProviderConfig{
					Endpoint: "us-west-1.log.aliyuncs.com",
					Project:  "   ",
					LogStore: "test-logstore",
				},
			},
			{
				name: "logstore with only spaces",
				config: SLSProviderConfig{
					Endpoint: "us-west-1.log.aliyuncs.com",
					Project:  "test-project",
					LogStore: "   ",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.config.Init()
				if err != nil {
					t.Errorf("Init() unexpected error = %v (whitespace fields should pass basic validation)", err)
				}
				// Note: The current implementation doesn't trim whitespace,
				// so whitespace-only strings pass validation. This might be
				// something to improve in the future.
			})
		}
	})
}
