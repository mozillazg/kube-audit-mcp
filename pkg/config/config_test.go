package config

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
)

// mockProvider implements the provider.Provider interface for testing
type mockProvider struct {
	name string
}

func (m *mockProvider) QueryAuditLog(_ context.Context, _ types.QueryAuditLogParams) (types.AuditLogResult, error) {
	return types.AuditLogResult{}, nil
}

func TestGetProviderByName(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		inputName      string
		expectedError  string
		expectedFound  bool
		setupProviders bool
	}{
		{
			name: "empty name uses default cluster",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "default-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "alibaba-sls",
							AlibabaSLS: &alibaba.SLSProviderConfig{
								Endpoint: "https://example.aliyuncs.com",
								Region:   "us-west-1",
								Project:  "test-project",
								LogStore: "test-logstore",
							},
						},
					},
				},
			},
			inputName:      "",
			expectedError:  "",
			expectedFound:  true,
			setupProviders: true,
		},
		{
			name: "find cluster by exact name match",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "test-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "alibaba-sls",
							AlibabaSLS: &alibaba.SLSProviderConfig{
								Endpoint: "https://example.aliyuncs.com",
								Region:   "us-west-1",
								Project:  "test-project",
								LogStore: "test-logstore",
							},
						},
					},
				},
			},
			inputName:      "test-cluster",
			expectedError:  "",
			expectedFound:  true,
			setupProviders: true,
		},
		{
			name: "find cluster by alias match",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "production-cluster",
						Alias:    []string{"prod", "production"},
						Disabled: false,
						Provider: ProviderConfig{
							Name: "aws-cloudwatch-logs",
							AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
								Region:             "us-east-1",
								LogGroupName:       "test-log-group",
								LogGroupIdentifier: "test-identifier",
							},
						},
					},
				},
			},
			inputName:      "prod",
			expectedError:  "",
			expectedFound:  true,
			setupProviders: true,
		},
		{
			name: "disabled cluster returns error",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "disabled-cluster",
						Disabled: true,
						Provider: ProviderConfig{
							Name: "alibaba-sls",
							AlibabaSLS: &alibaba.SLSProviderConfig{
								Endpoint: "https://example.aliyuncs.com",
								Region:   "us-west-1",
								Project:  "test-project",
								LogStore: "test-logstore",
							},
						},
					},
				},
			},
			inputName:     "disabled-cluster",
			expectedError: "cluster disabled-cluster is disabled",
			expectedFound: false,
		},
		{
			name: "cluster not found returns error",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "existing-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "alibaba-sls",
							AlibabaSLS: &alibaba.SLSProviderConfig{
								Endpoint: "https://example.aliyuncs.com",
								Region:   "us-west-1",
								Project:  "test-project",
								LogStore: "test-logstore",
							},
						},
					},
				},
			},
			inputName:     "non-existent-cluster",
			expectedError: "provider not found for name: non-existent-cluster",
			expectedFound: false,
		},
		{
			name: "invalid provider configuration returns error",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "invalid-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "unknown-provider",
						},
					},
				},
			},
			inputName:     "invalid-cluster",
			expectedError: "get provider for cluster invalid-cluster: unknown provider: unknown-provider",
			expectedFound: false,
		},
		{
			name: "missing required provider config returns error",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "missing-config-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "alibaba-sls",
							// Missing AlibabaSLS config
						},
					},
				},
			},
			inputName:     "missing-config-cluster",
			expectedError: "get provider for cluster missing-config-cluster: provider alibaba-sls requires alibaba_sls configuration",
			expectedFound: false,
		},
		{
			name: "aws provider with valid config",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "aws-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "aws-cloudwatch-logs",
							AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
								Region:             "us-east-1",
								LogGroupName:       "test-log-group",
								LogGroupIdentifier: "test-identifier",
							},
						},
					},
				},
			},
			inputName:      "aws-cluster",
			expectedError:  "",
			expectedFound:  true,
			setupProviders: true,
		},
		{
			name: "provider name with underscores normalized",
			config: &Config{
				DefaultCluster: "default-cluster",
				Clusters: []*Cluster{
					{
						Name:     "underscore-cluster",
						Disabled: false,
						Provider: ProviderConfig{
							Name: "aws_cloudwatch_logs", // Underscore version
							AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
								Region:             "us-east-1",
								LogGroupName:       "test-log-group",
								LogGroupIdentifier: "test-identifier",
							},
						},
					},
				},
			},
			inputName:      "underscore-cluster",
			expectedError:  "",
			expectedFound:  true,
			setupProviders: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup providers if needed (for cases where we expect success)
			if tt.setupProviders {
				for _, cluster := range tt.config.Clusters {
					if !cluster.Disabled {
						// Pre-create mock provider to avoid actual provider initialization
						cluster.p = &mockProvider{name: cluster.Name}
					}
				}
			}

			provider, err := tt.config.GetProviderByName(tt.inputName)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.expectedError)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
			}

			if tt.expectedFound {
				if provider == nil {
					t.Error("expected provider to be found, got nil")
				}
			} else {
				if provider != nil {
					t.Error("expected provider to be nil, got non-nil")
				}
			}
		})
	}
}

func TestGetProviderByName_ConcurrentAccess(t *testing.T) {
	config := &Config{
		DefaultCluster: "default-cluster",
		Clusters: []*Cluster{
			{
				Name:     "concurrent-cluster",
				Disabled: false,
				Provider: ProviderConfig{
					Name: "alibaba-sls",
					AlibabaSLS: &alibaba.SLSProviderConfig{
						Endpoint: "https://example.aliyuncs.com",
						Region:   "us-west-1",
						Project:  "test-project",
						LogStore: "test-logstore",
					},
				},
			},
		},
	}

	// Pre-setup the provider to avoid actual initialization
	config.Clusters[0].p = &mockProvider{name: "concurrent-cluster"}

	const numGoroutines = 10
	const numCalls = 100

	var wg sync.WaitGroup
	cerrors := make(chan error, numGoroutines*numCalls)

	// Start multiple goroutines making concurrent calls
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numCalls; j++ {
				provider, err := config.GetProviderByName("concurrent-cluster")
				if err != nil {
					cerrors <- err
					return
				}
				if provider == nil {
					cerrors <- errors.New("provider is nil")
					return
				}
			}
		}()
	}

	wg.Wait()
	close(cerrors)

	// Check for any errors
	for err := range cerrors {
		t.Errorf("concurrent access error: %v", err)
	}
}

func TestGetProviderByName_EmptyClustersList(t *testing.T) {
	config := &Config{
		DefaultCluster: "default-cluster",
		Clusters:       []*Cluster{}, // Empty clusters list
	}

	provider, err := config.GetProviderByName("any-name")
	if err == nil {
		t.Error("expected error for empty clusters list, got nil")
	}
	if provider != nil {
		t.Error("expected provider to be nil, got non-nil")
	}

	expectedError := "provider not found for name: any-name"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestGetProviderByName_MultipleAliasMatches(t *testing.T) {
	config := &Config{
		DefaultCluster: "default-cluster",
		Clusters: []*Cluster{
			{
				Name:     "cluster1",
				Alias:    []string{"prod", "production"},
				Disabled: false,
				Provider: ProviderConfig{
					Name: "alibaba-sls",
					AlibabaSLS: &alibaba.SLSProviderConfig{
						Endpoint: "https://example.aliyuncs.com",
						Region:   "us-west-1",
						Project:  "test-project",
						LogStore: "test-logstore",
					},
				},
			},
			{
				Name:     "cluster2",
				Alias:    []string{"prod", "staging"}, // Same alias "prod"
				Disabled: false,
				Provider: ProviderConfig{
					Name: "aws-cloudwatch-logs",
					AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
						Region:             "us-east-1",
						LogGroupName:       "test-log-group",
						LogGroupIdentifier: "test-identifier",
					},
				},
			},
		},
	}

	// Pre-setup providers
	config.Clusters[0].p = &mockProvider{name: "cluster1"}
	config.Clusters[1].p = &mockProvider{name: "cluster2"}

	// Should return the first matching cluster
	provider, err := config.GetProviderByName("prod")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Error("expected provider to be found, got nil")
	}

	// Verify it's the first cluster's provider
	mockProv, ok := provider.(*mockProvider)
	if !ok {
		t.Error("expected mock provider")
	} else if mockProv.name != "cluster1" {
		t.Errorf("expected first cluster provider, got provider for %s", mockProv.name)
	}
}

func TestCluster_createProvider(t *testing.T) {
	tests := []struct {
		name          string
		cluster       *Cluster
		expectedError string
		expectedType  string
	}{
		{
			name: "create alibaba sls provider successfully",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "alibaba-sls",
					AlibabaSLS: &alibaba.SLSProviderConfig{
						Endpoint: "https://test.aliyuncs.com",
						Region:   "us-west-1",
						Project:  "test-project",
						LogStore: "test-logstore",
					},
				},
			},
			expectedError: "",
			expectedType:  "*alibaba.SLSProvider",
		},
		{
			name: "create alibaba sls provider with normalized name",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "alibaba_sls", // underscore should be normalized to dash
					AlibabaSLS: &alibaba.SLSProviderConfig{
						Endpoint: "https://test.aliyuncs.com",
						Region:   "us-west-1",
						Project:  "test-project",
						LogStore: "test-logstore",
					},
				},
			},
			expectedError: "",
			expectedType:  "*alibaba.SLSProvider",
		},
		{
			name: "create aws cloudwatch logs provider successfully",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "aws-cloudwatch-logs",
					AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
						Region:       "us-east-1",
						LogGroupName: "test-log-group",
					},
				},
			},
			expectedError: "",
			expectedType:  "*aws.CloudWatchLogsProvider",
		},
		{
			name: "create aws cloudwatch logs provider with normalized name",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "aws_cloudwatch_logs", // underscore should be normalized to dash
					AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
						Region:             "us-east-1",
						LogGroupIdentifier: "test-identifier",
					},
				},
			},
			expectedError: "",
			expectedType:  "*aws.CloudWatchLogsProvider",
		},
		{
			name: "alibaba sls provider missing configuration",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name:       "alibaba-sls",
					AlibabaSLS: nil,
				},
			},
			expectedError: "provider alibaba-sls requires alibaba_sls configuration",
			expectedType:  "",
		},
		{
			name: "aws cloudwatch logs provider missing configuration",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name:              "aws-cloudwatch-logs",
					AwsCloudWatchLogs: nil,
				},
			},
			expectedError: "provider aws-cloudwatch-logs requires aws_cloudwatch_logs configuration",
			expectedType:  "",
		},
		{
			name: "unknown provider",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "unknown-provider",
				},
			},
			expectedError: "unknown provider: unknown-provider",
			expectedType:  "",
		},
		{
			name: "alibaba sls provider with invalid configuration",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "alibaba-sls",
					AlibabaSLS: &alibaba.SLSProviderConfig{
						// Missing required fields (endpoint, project, logstore)
						Region: "us-west-1",
					},
				},
			},
			expectedError: "init provider alibaba-sls:",
			expectedType:  "",
		},
		{
			name: "aws cloudwatch logs provider with invalid configuration",
			cluster: &Cluster{
				Name: "test-cluster",
				Provider: ProviderConfig{
					Name: "aws-cloudwatch-logs",
					AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
						// Missing required fields (log_group_name)
						Region: "us-east-1",
					},
				},
			},
			expectedError: "init provider aws-cloudwatch-logs:",
			expectedType:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := tt.cluster.createProvider()

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error containing '%s', but got no error", tt.expectedError)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', but got '%s'", tt.expectedError, err.Error())
				}
				if provider != nil {
					t.Error("expected provider to be nil when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got: %v", err)
				}
				if provider == nil {
					t.Error("expected provider to be non-nil when no error occurs")
				} else {
					actualType := fmt.Sprintf("%T", provider)
					if actualType != tt.expectedType {
						t.Errorf("expected provider type '%s', but got '%s'", tt.expectedType, actualType)
					}
				}
			}
		})
	}
}

func TestProviderConfig_normalizedName(t *testing.T) {
	tests := []struct {
		name     string
		provider ProviderConfig
		expected string
	}{
		{
			name: "name with underscores",
			provider: ProviderConfig{
				Name: "alibaba_sls",
			},
			expected: "alibaba-sls",
		},
		{
			name: "name with multiple underscores",
			provider: ProviderConfig{
				Name: "aws_cloudwatch_logs",
			},
			expected: "aws-cloudwatch-logs",
		},
		{
			name: "name with dashes (no change)",
			provider: ProviderConfig{
				Name: "alibaba-sls",
			},
			expected: "alibaba-sls",
		},
		{
			name: "name without underscores or dashes",
			provider: ProviderConfig{
				Name: "simpleProvider",
			},
			expected: "simpleProvider",
		},
		{
			name: "empty name",
			provider: ProviderConfig{
				Name: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.provider.normalizedName()
			if result != tt.expected {
				t.Errorf("expected '%s', but got '%s'", tt.expected, result)
			}
		})
	}
}
