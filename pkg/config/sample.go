package config

import (
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/gcp"
)

var SampleConfig = Config{
	DefaultCluster: "dev",
	Clusters: []*Cluster{
		{
			Name:     "prod",
			Alias:    []string{"aws-prod"},
			Disabled: false,
			Provider: ProviderConfig{
				Name: aws.CloudWatchProviderName,
				AwsCloudWatchLogs: &aws.CloudWatchLogsProviderConfig{
					Region:             "",
					LogGroupName:       "/aws/eks/xxx/cluster",
					LogGroupIdentifier: "",
				},
			},
		},
		{
			Name:     "dev",
			Alias:    []string{"dev-cluster"},
			Disabled: false,
			Provider: ProviderConfig{
				Name: alibaba.SLSProviderName,
				AlibabaSLS: &alibaba.SLSProviderConfig{
					Endpoint: "cn-hangzhou.log.aliyuncs.com",
					Project:  "k8s-cxxx",
					LogStore: "audit-cxxx",
				},
			},
		},
		{
			Name:        "test",
			Description: "",
			Alias:       nil,
			Disabled:    false,
			Provider: ProviderConfig{
				Name: gcp.CloudLoggingProviderName,
				GcpCloudLogging: &gcp.CloudLoggingProviderConfig{
					ProjectId: "test-233xxx",
				},
			},
		},
	},
}
