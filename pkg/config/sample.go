package config

import (
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
)

var SampleConfig = Config{
	DefaultCluster: "dev",
	Clusters: []*Cluster{
		{
			Name:     "prod",
			Alias:    []string{"eks", "aws-cluster", "aws_eks", "aws-prod"},
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
			Alias:    []string{"cxxx", "dev-cluster"},
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
	},
}
