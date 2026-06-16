package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "ponte",
		Short: "Sync instructions and skills across AI agent vendors",
	}
	root.AddCommand(newSyncCommand())
	root.AddCommand(newSyspromptCommand())
	root.AddCommand(newManualCommand())
	return root
}
