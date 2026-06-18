package sync

import (
	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	"github.com/flexksx/ponte/apps/ponte/internal/config"
	"github.com/flexksx/ponte/apps/ponte/internal/skill"
	"github.com/flexksx/ponte/apps/ponte/internal/store"
	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
)

type UseCase struct {
	ReadSystemPrompt      systemprompt.Reader
	ReadConfig            config.ConfigReader
	GetAgentConfiguration agentvendor.ConfigurationPort
	ResolveSkill          skill.Resolver
	BuildGeneration       store.GenerationBuilder
	ComputeHash           store.HashComputer
	ActivateForVendor     store.VendorActivator
}

type resolvedVendor struct {
	name   agentvendor.AgentVendorName
	config agentvendor.AgentVendorConfiguration
}

func (u *UseCase) Execute(request SyncRequest) (SyncResult, error) {
	targetNames, err := u.resolveTargets(request.TargetAgents)
	if err != nil {
		return SyncResult{}, err
	}

	vendors, err := u.resolveVendors(targetNames)
	if err != nil {
		return SyncResult{}, err
	}

	prompt, err := u.resolveSystemPrompt(request.SystemPromptOverride)
	if err != nil {
		return SyncResult{}, err
	}

	input, err := ResolveBuildInput(prompt, request.Skills, request.Subagents, u.ResolveSkill)
	if err != nil {
		return SyncResult{}, err
	}

	if request.DryRun {
		hash, err := u.ComputeHash(input)
		if err != nil {
			return SyncResult{}, err
		}
		return SyncResult{GenerationHash: hash, Targets: targetNames, DryRun: true}, nil
	}

	generation, err := u.BuildGeneration(input)
	if err != nil {
		return SyncResult{}, err
	}

	for _, vendor := range vendors {
		if err := u.ActivateForVendor(generation, vendor.config.GlobalInstructionFilePath, vendor.config.SkillsDirectoryPath, vendor.config.SubagentsDirectoryPath); err != nil {
			return SyncResult{}, err
		}
	}
	return SyncResult{GenerationHash: generation.Hash, Targets: targetNames, DryRun: false}, nil
}

func (u *UseCase) resolveTargets(requested []agentvendor.AgentVendorName) ([]agentvendor.AgentVendorName, error) {
	if len(requested) > 0 {
		return requested, nil
	}
	cfg, err := u.ReadConfig()
	if err != nil {
		return nil, err
	}
	var enabled []agentvendor.AgentVendorName
	for name, entry := range cfg.Agents {
		if entry.Enabled {
			enabled = append(enabled, name)
		}
	}
	if len(enabled) == 0 {
		return nil, ErrNoAgentsConfigured{}
	}
	return enabled, nil
}

// resolveVendors validates every target up front so an unknown agent fails
// before anything is built — true for a dry run as much as a real sync.
func (u *UseCase) resolveVendors(names []agentvendor.AgentVendorName) ([]resolvedVendor, error) {
	vendors := make([]resolvedVendor, 0, len(names))
	for _, name := range names {
		vendorConfig, err := u.GetAgentConfiguration(name)
		if err != nil {
			return nil, ErrUnknownAgent{Name: name}
		}
		vendors = append(vendors, resolvedVendor{name: name, config: vendorConfig})
	}
	return vendors, nil
}

func (u *UseCase) resolveSystemPrompt(override *systemprompt.SystemPrompt) (systemprompt.SystemPrompt, error) {
	if override != nil {
		return *override, nil
	}
	return u.ReadSystemPrompt()
}
