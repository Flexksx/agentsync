package sync

import (
	"fmt"

	"github.com/flexksx/ponte/apps/ponte/internal/config"
	"github.com/flexksx/ponte/apps/ponte/internal/skill"
	"github.com/flexksx/ponte/apps/ponte/internal/store"
	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
)

// ResolveBuildInput turns a system prompt plus the declared skill and subagent
// sources into the store.BuildInput a generation is built from. Both the sync
// use case and the status preview share it so a previewed generation hash
// always matches what a real sync produces.
func ResolveBuildInput(
	prompt systemprompt.SystemPrompt,
	skills []config.SkillEntry,
	subagents []config.SubagentEntry,
	resolveSkill skill.Resolver,
) (store.BuildInput, error) {
	resolvedSkills, err := resolveSkills(skills, resolveSkill)
	if err != nil {
		return store.BuildInput{}, err
	}
	resolvedSubagents, err := resolveSubagents(subagents, resolveSkill)
	if err != nil {
		return store.BuildInput{}, err
	}
	return store.BuildInput{
		SystemPromptContent: prompt.Content,
		Skills:              resolvedSkills,
		Subagents:           resolvedSubagents,
	}, nil
}

func resolveSkills(entries []config.SkillEntry, resolveSkill skill.Resolver) ([]store.ResolvedSkill, error) {
	resolved := make([]store.ResolvedSkill, 0, len(entries))
	for _, entry := range entries {
		dir, err := resolveSkill(entry.Source)
		if err != nil {
			return nil, fmt.Errorf("resolving skill %q: %w", entry.Name, err)
		}
		resolved = append(resolved, store.ResolvedSkill{Name: entry.Name, SourceDir: dir})
	}
	return resolved, nil
}

func resolveSubagents(entries []config.SubagentEntry, resolveSkill skill.Resolver) ([]store.ResolvedSubagent, error) {
	resolved := make([]store.ResolvedSubagent, 0, len(entries))
	for _, entry := range entries {
		dir, err := resolveSkill(entry.Source)
		if err != nil {
			return nil, fmt.Errorf("resolving subagent %q: %w", entry.Name, err)
		}
		resolved = append(resolved, store.ResolvedSubagent{Name: entry.Name, SourceDir: dir})
	}
	return resolved, nil
}
