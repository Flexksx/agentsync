package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flexksx/ponte/apps/ponte/internal/skill"
)

const (
	entryNameHeader = "NAME"
	entryTypeHeader = "TYPE"
	entryTypeWidth  = 5
)

// configEntry is a name plus its source, the shape both skills and subagents
// share. Rendering them through one path keeps the two listing commands
// identical in format.
type configEntry struct {
	name   string
	source skill.SkillSource
}

// printConfigEntries renders a name/type/source table, or an empty-state line
// keyed by noun ("skills", "subagents").
func printConfigEntries(cmd *cobra.Command, noun string, entries []configEntry) {
	out := cmd.OutOrStdout()
	if len(entries) == 0 {
		_, _ = fmt.Fprintf(out, "No %s configured.\n", noun)
		return
	}

	nameWidth := len(entryNameHeader)
	for _, entry := range entries {
		if len(entry.name) > nameWidth {
			nameWidth = len(entry.name)
		}
	}

	_, _ = fmt.Fprintf(out, "%-*s  %-*s  %s\n", nameWidth, entryNameHeader, entryTypeWidth, entryTypeHeader, "SOURCE")
	for _, entry := range entries {
		_, _ = fmt.Fprintf(out, "%-*s  %-*s  %s\n", nameWidth, entry.name, entryTypeWidth, entry.source.Type, formatSkillSource(entry.source))
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
