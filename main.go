package main

import (
	"log"

	"github.com/mozillazg/kube-audit-mcp/pkg/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
