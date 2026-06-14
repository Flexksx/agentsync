package main

import (
	"os"

	"github.com/flexksx/agentsync/apps/agentsync/internal/cli"
)

func main() {
	root := cli.NewRootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
