package cli

import (
	"fmt"
	"log"
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
	if s, err := os.Stat(savePath); err == nil {
		if !s.IsDir() {
			log.Fatalf("File %s already exists. Skipping save.", savePath)
			return
		}
	} else if !os.IsNotExist(err) {
		log.Fatalf("Failed to check file %s: %+v", savePath, err)
		return
	}

	dirPath := path.Dir(savePath)
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		log.Fatalf("Failed to create directory %s: %+v", dirPath, err)
		return
	}
	if err := os.WriteFile(savePath, sampleConf, 0600); err != nil {
		log.Fatalf("Failed to write sample config to %s: %+v", savePath, err)
		return
	}

	log.Printf("Sample configuration saved to %s", savePath)
}
