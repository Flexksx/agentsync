package cli

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed docs/manual.md
var manualContent string

func newManualCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "manual",
		Short: "Show the full configuration and usage guide",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(manualContent)
		},
	}
}
