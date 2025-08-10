package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/mozillazg/kube-audit-mcp/pkg/config"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var saveSampleConf bool

var sampleConfCmd = &cobra.Command{
	Use:     "sample-config",
	Aliases: []string{"sample-conf", "sample_conf", "sample_config"},
	Short:   "Print a sample configuration for kube-audit-mcp.",
	Run: func(cmd *cobra.Command, args []string) {
		runSampleConfCmd(cmd, args)
	},
}

func init() {
	sampleConfCmd.Flags().BoolVarP(
		&saveSampleConf, "save", "s",
		false,
		fmt.Sprintf("Save the sample configuration to file %s.",
			config.DefaultConfigFile()))
}

func runSampleConfCmd(cmd *cobra.Command, args []string) {
	sampleConf, _ := yaml.Marshal(config.SampleConfig)
	fmt.Println()
	fmt.Println(string(sampleConf))

	if !saveSampleConf {
		return
	}

	savePath := config.DefaultConfigFile()
	dirPath := path.Dir(savePath)
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dirPath, err)
		return
	}
	if err := os.WriteFile(savePath, sampleConf, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write sample config to %s: %v\n", savePath, err)
		return
	}
	fmt.Fprintf(os.Stderr, "Sample configuration saved to %s\n", savePath)
}
