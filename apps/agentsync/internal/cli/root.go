package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "agentsync",
		Short: "Sync instructions and skills across AI agent vendors",
	}
	root.AddCommand(newSyncCommand())
	root.AddCommand(newSyspromptCommand())
	return root
}
