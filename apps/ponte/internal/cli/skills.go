package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flexksx/ponte/apps/ponte/internal/config"
	configadapter "github.com/flexksx/ponte/apps/ponte/internal/config/adapter"
	"github.com/flexksx/ponte/apps/ponte/internal/skill"
)

const (
	skillNameHeader = "NAME"
	skillTypeHeader = "TYPE"
	skillTypeWidth  = 5
)

func newSkillsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "skills",
		Short: "List the skills declared in config.toml",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := configadapter.ReadConfig()
			if err != nil {
				return err
			}
			printSkillEntries(cmd, cfg.Skills)
			return nil
		},
	}
}

func printSkillEntries(cmd *cobra.Command, skills []config.SkillEntry) {
	out := cmd.OutOrStdout()
	if len(skills) == 0 {
		_, _ = fmt.Fprintln(out, "No skills configured.")
		return
	}

	nameWidth := len(skillNameHeader)
	for _, entry := range skills {
		if len(entry.Name) > nameWidth {
			nameWidth = len(entry.Name)
		}
	}

	_, _ = fmt.Fprintf(out, "%-*s  %-*s  %s\n", nameWidth, skillNameHeader, skillTypeWidth, skillTypeHeader, "SOURCE")
	for _, entry := range skills {
		_, _ = fmt.Fprintf(out, "%-*s  %-*s  %s\n", nameWidth, entry.Name, skillTypeWidth, entry.Source.Type, formatSkillSource(entry.Source))
	}
}

func formatSkillSource(source skill.SkillSource) string {
	switch source.Type {
	case skill.LocalSourceType:
		return source.LocalPath
	case skill.GitSourceType:
		description := source.GitURL
		if source.GitRef != "" {
			description += "@" + source.GitRef
		}
		if source.Subdir != "" {
			description += " (subdir: " + source.Subdir + ")"
		}
		return description
	default:
		return string(source.Type)
	}
}
