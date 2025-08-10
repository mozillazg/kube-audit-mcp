package config

import (
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
)

var SampleConfig = Config{
	ProviderName: aws.CloudWatchProviderName,
	ProviderConfig: ProviderConfig{
		AlibabaSLS: alibaba.SLSProviderConfig{
			Endpoint: "cn-hangzhou.log.aliyuncs.com",
			Project:  "k8s-cxxx",
			LogStore: "audit-cxxx",
		},
		AwsCloudWatchLogs: aws.CloudWatchLogsProviderConfig{
			Region:             "",
			LogGroupName:       "/aws/eks/xxx/cluster",
			LogGroupIdentifier: "",
		},
	},
}
