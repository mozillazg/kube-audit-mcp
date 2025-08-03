package cli

import (
	"fmt"

	"github.com/mozillazg/kube-audit-mcp/pkg/config"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var sampleConfCmd = &cobra.Command{
	Use:     "sample-config",
	Aliases: []string{"sample-conf", "sample_conf", "sample_config"},
	Short:   "Print a sample configuration for kube-audit-mcp.",
	Run: func(cmd *cobra.Command, args []string) {
		sampleConf, _ := yaml.Marshal(config.SampleConfig)
		fmt.Println(string(sampleConf))
	},
}
