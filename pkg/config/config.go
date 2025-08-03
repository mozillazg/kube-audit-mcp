package config

import (
	"fmt"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"os"

	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"sigs.k8s.io/yaml"
)

type Config struct {
	ProviderName   string         `yaml:"provider_name" json:"provider_name"`
	ProviderConfig ProviderConfig `yaml:"provider_config" json:"provider_config"`
}

type ProviderConfig struct {
	AlibabaSLS alibaba.SLSProviderConfig `yaml:"alibaba_sls" json:"alibaba_sls"`
}

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
	},
}

func NewConfigFromFile(filePath string) (*Config, error) {
	config := &Config{}
	if err := config.LoadFromFile(filePath); err != nil {
		return nil, fmt.Errorf("load config from file: %w", err)
	}
	return config, nil
}

func (c *Config) LoadFromFile(filePath string) error {
	// Implement the logic to load the configuration from a YAML file.
	// This is a placeholder for the actual implementation.
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}
	return yaml.Unmarshal(data, c)
}

func (c *Config) NewProvider() (provider.Provider, error) {
	switch c.ProviderName {
	case alibaba.SLSProviderName:
		p, err := alibaba.NewSLSProvider(&c.ProviderConfig.AlibabaSLS)
		if err != nil {
			return nil, fmt.Errorf("init provider %s: %w", c.ProviderName, err)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", c.ProviderName)
	}
}
