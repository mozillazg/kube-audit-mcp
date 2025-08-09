package config

import (
	"fmt"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
	"os"
	"strings"

	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"sigs.k8s.io/yaml"
)

type Config struct {
	ProviderName   string         `yaml:"provider_name" json:"provider_name"`
	ProviderConfig ProviderConfig `yaml:"provider_config" json:"provider_config"`
}

type ProviderConfig struct {
	AlibabaSLS        alibaba.SLSProviderConfig        `yaml:"alibaba_sls" json:"alibaba_sls"`
	AwsCloudWatchLogs aws.CloudWatchLogsProviderConfig `yaml:"aws_cloudwatch_logs" json:"aws_cloudwatch_logs"`
}

func NewConfigFromFile(filePath string) (*Config, error) {
	config := &Config{}
	if err := config.LoadFromFile(filePath); err != nil {
		return nil, fmt.Errorf("load config from file: %w", err)
	}
	return config, nil
}

func (c *Config) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}
	return yaml.Unmarshal(data, c)
}

func (c *Config) NewProvider() (provider.Provider, error) {
	switch strings.Replace(c.ProviderName, "_", "-", -1) {
	case alibaba.SLSProviderName:
		p, err := alibaba.NewSLSProvider(&c.ProviderConfig.AlibabaSLS)
		if err != nil {
			return nil, fmt.Errorf("init provider %s: %w", c.ProviderName, err)
		}
		return p, nil
	case aws.CloudWatchProviderName:
		p, err := aws.NewCloudWatchLogsProvider(&c.ProviderConfig.AwsCloudWatchLogs)
		if err != nil {
			return nil, fmt.Errorf("init provider %s: %w", c.ProviderName, err)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", c.ProviderName)
	}
}
