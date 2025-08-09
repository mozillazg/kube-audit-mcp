package config

import (
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
)

var SampleConfig = Config{
	ProviderName: alibaba.SLSProviderName,
	ProviderConfig: ProviderConfig{
		AlibabaSLS: alibaba.SLSProviderConfig{
			Endpoint:    "cn-hangzhou.log.aliyuncs.com",
			Region:      "cn-hangzhou",
			AuthVersion: "v4",
			Project:     "k8s-cxxx",
			LogStore:    "audit-cxxx",
		},
		AwsCloudWatchLogs: aws.CloudWatchLogsProviderConfig{
			Region:             "",
			LogGroupName:       "/aws/eks/xxx/cluster",
			LogGroupIdentifier: "",
		},
	},
}
