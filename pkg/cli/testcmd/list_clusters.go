package testcmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var listClustersmd = &cobra.Command{
	Use:     "list-clusters",
	Aliases: []string{"list_clusters", "list-clusters", "listclusters"},
	Short:   "call list_clusters",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		err := runListClustersCmd(ctx, os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	testCmd.AddCommand(listClustersmd)
}

func runListClustersCmd(ctx context.Context, cmd string) error {
	args := []string{"mcp"}
	return callTool(ctx, cmd, args, "list_clusters", nil)
}
