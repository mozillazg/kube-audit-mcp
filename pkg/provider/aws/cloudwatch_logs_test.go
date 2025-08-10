package aws

import (
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	"testing"
)

func TestCloudWatchLogsProvider_buildQuery(t *testing.T) {
	provider := &CloudWatchLogsProvider{}

	tests := []struct {
		name     string
		params   types.QueryAuditLogParams
		expected string
	}{
		{
			name: "basic query with only limit",
			params: types.QueryAuditLogParams{
				Limit: 100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with user exact match",
			params: types.QueryAuditLogParams{
				User:  "john.doe",
				Limit: 50,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "john.doe" | sort @timestamp desc | limit 50`,
		},
		{
			name: "query with user wildcard",
			params: types.QueryAuditLogParams{
				User:  "admin*",
				Limit: 25,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username like "admin." | sort @timestamp desc | limit 25`,
		},
		{
			name: "query with user asterisk (should be ignored)",
			params: types.QueryAuditLogParams{
				User:  "*",
				Limit: 100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with namespace exact match",
			params: types.QueryAuditLogParams{
				Namespace: "default",
				Limit:     100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.namespace = "default" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with namespace wildcard",
			params: types.QueryAuditLogParams{
				Namespace: "kube-*",
				Limit:     100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.namespace like "kube-." | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with namespace asterisk (should be ignored)",
			params: types.QueryAuditLogParams{
				Namespace: "*",
				Limit:     100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with single verb",
			params: types.QueryAuditLogParams{
				Verbs: []string{"get"},
				Limit: 100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter verb in ["get"] | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with multiple verbs",
			params: types.QueryAuditLogParams{
				Verbs: []string{"get", "create", "update"},
				Limit: 100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter verb in ["get", "create", "update"] | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with single resource type",
			params: types.QueryAuditLogParams{
				ResourceTypes: []string{"pods"},
				Limit:         100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.resource in ["pods"] | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with multiple resource types",
			params: types.QueryAuditLogParams{
				ResourceTypes: []string{"pods", "services", "deployments"},
				Limit:         100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.resource in ["pods", "services", "deployments"] | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with resource name exact match",
			params: types.QueryAuditLogParams{
				ResourceName: "my-pod",
				Limit:        100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.name = "my-pod" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with resource name wildcard",
			params: types.QueryAuditLogParams{
				ResourceName: "app-*",
				Limit:        100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.name like "app-." | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with resource name asterisk (should be ignored)",
			params: types.QueryAuditLogParams{
				ResourceName: "*",
				Limit:        100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with all parameters (exact matches)",
			params: types.QueryAuditLogParams{
				User:          "john.doe",
				Namespace:     "production",
				Verbs:         []string{"get", "list"},
				ResourceTypes: []string{"pods", "services"},
				ResourceName:  "web-server",
				Limit:         200,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "john.doe" and objectRef.namespace = "production" and verb in ["get", "list"] and objectRef.resource in ["pods", "services"] and objectRef.name = "web-server" | sort @timestamp desc | limit 200`,
		},
		{
			name: "query with all parameters (wildcards)",
			params: types.QueryAuditLogParams{
				User:          "admin*",
				Namespace:     "kube-*",
				Verbs:         []string{"create", "update", "delete"},
				ResourceTypes: []string{"configmaps", "secrets"},
				ResourceName:  "app-*",
				Limit:         150,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username like "admin." and objectRef.namespace like "kube-." and verb in ["create", "update", "delete"] and objectRef.resource in ["configmaps", "secrets"] and objectRef.name like "app-." | sort @timestamp desc | limit 150`,
		},
		{
			name: "query with mixed exact and wildcard parameters",
			params: types.QueryAuditLogParams{
				User:          "service*",
				Namespace:     "default",
				Verbs:         []string{"patch"},
				ResourceTypes: []string{"deployments"},
				Limit:         75,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username like "service." and objectRef.namespace = "default" and verb in ["patch"] and objectRef.resource in ["deployments"] | sort @timestamp desc | limit 75`,
		},
		{
			name: "query with empty user string",
			params: types.QueryAuditLogParams{
				User:      "",
				Namespace: "test",
				Limit:     100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter objectRef.namespace = "test" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with empty namespace string",
			params: types.QueryAuditLogParams{
				User:      "testuser",
				Namespace: "",
				Limit:     100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "testuser" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with empty resource name string",
			params: types.QueryAuditLogParams{
				User:         "testuser",
				ResourceName: "",
				Limit:        100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "testuser" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with empty verbs slice",
			params: types.QueryAuditLogParams{
				User:  "testuser",
				Verbs: []string{},
				Limit: 100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "testuser" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with empty resource types slice",
			params: types.QueryAuditLogParams{
				User:          "testuser",
				ResourceTypes: []string{},
				Limit:         100,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "testuser" | sort @timestamp desc | limit 100`,
		},
		{
			name: "query with zero limit",
			params: types.QueryAuditLogParams{
				User:  "testuser",
				Limit: 0,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "testuser" | sort @timestamp desc | limit 0`,
		},
		{
			name: "query with large limit",
			params: types.QueryAuditLogParams{
				User:  "testuser",
				Limit: 10000,
			},
			expected: `fields @message | filter @logStream like "kube-apiserver-audit" | filter user.username = "testuser" | sort @timestamp desc | limit 10000`,
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

func TestGetFilterExp(t *testing.T) {
	tests := []struct {
		name        string
		keyword     string
		expectedExp string
		expectedVal string
	}{
		{
			name:        "exact match - simple string",
			keyword:     "test",
			expectedExp: "=",
			expectedVal: "test",
		},
		{
			name:        "exact match - with spaces",
			keyword:     "test user",
			expectedExp: "=",
			expectedVal: "test user",
		},
		{
			name:        "exact match - with special characters",
			keyword:     "test@example.com",
			expectedExp: "=",
			expectedVal: "test@example.com",
		},
		{
			name:        "wildcard match - simple",
			keyword:     "test*",
			expectedExp: "like",
			expectedVal: "test.",
		},
		{
			name:        "wildcard match - with prefix",
			keyword:     "admin*",
			expectedExp: "like",
			expectedVal: "admin.",
		},
		{
			name:        "wildcard match - empty prefix",
			keyword:     "*",
			expectedExp: "like",
			expectedVal: ".",
		},
		{
			name:        "wildcard match - with special characters",
			keyword:     "user@domain*",
			expectedExp: "like",
			expectedVal: "user@domain.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, val := getFilterExp(tt.keyword)
			if exp != tt.expectedExp {
				t.Errorf("getFilterExp() exp = %q, want %q", exp, tt.expectedExp)
			}
			if val != tt.expectedVal {
				t.Errorf("getFilterExp() val = %q, want %q", val, tt.expectedVal)
			}
		})
	}
}

func TestCloudWatchLogsProviderConfig_Init(t *testing.T) {
	tests := []struct {
		name    string
		config  CloudWatchLogsProviderConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with log_group_name",
			config: CloudWatchLogsProviderConfig{
				LogGroupName: "test-group",
			},
			wantErr: false,
		},
		{
			name: "valid config with log_group_identifier",
			config: CloudWatchLogsProviderConfig{
				LogGroupIdentifier: "test-identifier",
			},
			wantErr: false,
		},
		{
			name: "invalid config - both log_group_name and identifier provided",
			config: CloudWatchLogsProviderConfig{
				LogGroupName:       "test-group",
				LogGroupIdentifier: "test-identifier",
			},
			wantErr: true,
			errMsg:  "only one of log_group_name or log_group_identifier can be provided",
		},
		{
			name:    "invalid config - neither log_group_name nor identifier provided",
			config:  CloudWatchLogsProviderConfig{},
			wantErr: true,
			errMsg:  "either log_group_name or log_group_identifier must be provided",
		},
		{
			name: "valid config with region and log_group_name",
			config: CloudWatchLogsProviderConfig{
				Region:       "us-west-2",
				LogGroupName: "test-group",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Init()
			if tt.wantErr {
				if err == nil {
					t.Error("Init() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Init() error = %v, want %v", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Init() unexpected error = %v", err)
			}
		})
	}
}
