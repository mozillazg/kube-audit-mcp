package testcmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var listCommonResourceTypesmd = &cobra.Command{
	Use:     "list-common-resource-types",
	Aliases: []string{"list_common_resource_types", "list-common-resource-types", "listcommonresourcetypes"},
	Short:   "call list_common_resource_types",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		err := runListCommonResourceTypesCmd(ctx, os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	testCmd.AddCommand(listCommonResourceTypesmd)
}

func runListCommonResourceTypesCmd(ctx context.Context, cmd string) error {
	args := []string{"mcp"}
	return callTool(ctx, cmd, args, "list_common_resource_types", nil)
}
