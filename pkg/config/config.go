package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/alibaba"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider/aws"
	"github.com/mozillazg/kube-audit-mcp/pkg/utils"
	"sigs.k8s.io/yaml"
)

type Config struct {
	DefaultCluster string     `yaml:"default_cluster" json:"default_cluster"`
	Clusters       []*Cluster `yaml:"clusters,omitempty" json:"clusters,omitempty"`
	mu             sync.RWMutex
}

type Cluster struct {
	Name     string   `yaml:"name" json:"name"`
	Alias    []string `yaml:"alias,omitempty" json:"alias,omitempty"`
	Disabled bool     `yaml:"disabled" json:"disabled"`

	Provider ProviderConfig `yaml:"provider" json:"provider"`

	p  provider.Provider
	mu sync.RWMutex
}

type ProviderConfig struct {
	Name              string                            `yaml:"name" json:"name"`
	AlibabaSLS        *alibaba.SLSProviderConfig        `yaml:"alibaba_sls,omitempty" json:"alibaba_sls,omitempty"`
	AwsCloudWatchLogs *aws.CloudWatchLogsProviderConfig `yaml:"aws_cloudwatch_logs,omitempty" json:"aws_cloudwatch_logs,omitempty"`
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

func (c *Config) Init() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var clusterNames []string

	for _, cluster := range c.Clusters {
		if cluster.Disabled {
			continue // Skip disabled clusters
		}
		if cluster.Name == "" {
			return errors.New("cluster name is required")
		}
		if cluster.Provider.Name == "" {
			return fmt.Errorf("provider name is required for cluster %s", cluster.Name)
		}

		// Initialize provider for validation
		_, err := cluster.getProvider()
		if err != nil {
			return fmt.Errorf("get provider for cluster %s: %w", cluster.Name, err)
		}

		clusterNames = append(clusterNames, cluster.Name)
		clusterNames = append(clusterNames, cluster.Alias...)
	}

	if c.DefaultCluster == "" {
		return errors.New("default_cluster is required")
	}
	if !utils.Contains(clusterNames, c.DefaultCluster) {
		return fmt.Errorf("default_cluster %s not found in clusters", c.DefaultCluster)
	}

	return nil
}

func (c *Config) GetProviderByName(name string) (provider.Provider, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if name == "" {
		name = c.DefaultCluster
	}
	for _, cluster := range c.Clusters {
		if cluster.Name == name || utils.Contains(cluster.Alias, name) {
			// Check if cluster is disabled
			if cluster.Disabled {
				return nil, fmt.Errorf("cluster %s is disabled", cluster.Name)
			}

			p, err := cluster.getProvider()
			if err != nil {
				return nil, fmt.Errorf("get provider for cluster %s: %w", cluster.Name, err)
			}
			return p, nil
		}
	}

	return nil, fmt.Errorf("provider not found for name: %s", name)
}

func (c *Config) AvailableClusterNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var names []string
	for _, cluster := range c.Clusters {
		if cluster.Disabled {
			continue
		}
		names = append(names, cluster.Name)
		names = append(names, cluster.Alias...)
	}
	names = utils.RemoveDuplicates(names)
	return names
}

func (p ProviderConfig) normalizedName() string {
	return strings.Replace(p.Name, "_", "-", -1)
}

// createProvider creates a new provider instance based on the provider configuration
func (c *Cluster) createProvider() (provider.Provider, error) {
	pconfig := c.Provider
	normalizedName := pconfig.normalizedName()

	switch normalizedName {
	case alibaba.SLSProviderName:
		if pconfig.AlibabaSLS == nil {
			return nil, fmt.Errorf("provider %s requires alibaba_sls configuration", pconfig.Name)
		}
		p, err := alibaba.NewSLSProvider(pconfig.AlibabaSLS)
		if err != nil {
			return nil, fmt.Errorf("init provider %s: %w", pconfig.Name, err)
		}
		return p, nil
	case aws.CloudWatchProviderName:
		if pconfig.AwsCloudWatchLogs == nil {
			return nil, fmt.Errorf("provider %s requires aws_cloudwatch_logs configuration", pconfig.Name)
		}
		p, err := aws.NewCloudWatchLogsProvider(pconfig.AwsCloudWatchLogs)
		if err != nil {
			return nil, fmt.Errorf("init provider %s: %w", pconfig.Name, err)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", pconfig.Name)
	}
}

func (c *Cluster) getProvider() (provider.Provider, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.p != nil {
		return c.p, nil
	}

	if c.Disabled {
		return nil, fmt.Errorf("cluster %s is disabled", c.Name)
	}

	p, err := c.createProvider()
	if err != nil {
		return nil, err
	}

	c.p = p
	return p, nil
}
