package cli

import (
	"github.com/spf13/cobra"

	configadapter "github.com/flexksx/ponte/apps/ponte/internal/config/adapter"
)

func newSubagentsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "subagents",
		Short: "List the subagents declared in config.toml",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := configadapter.ReadConfig()
			if err != nil {
				return err
			}
			entries := make([]configEntry, 0, len(cfg.Subagents))
			for _, subagentEntry := range cfg.Subagents {
				entries = append(entries, configEntry{name: subagentEntry.Name, source: subagentEntry.Source})
			}
			printConfigEntries(cmd, "subagents", entries)
			return nil
		},
	}
}
