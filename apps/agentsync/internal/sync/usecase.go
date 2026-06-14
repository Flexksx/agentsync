package sync

import (
	"github.com/flexksx/agentsync/apps/agentsync/internal/agentvendor"
	"github.com/flexksx/agentsync/apps/agentsync/internal/config"
	"github.com/flexksx/agentsync/apps/agentsync/internal/systemprompt"
)

type UseCase struct {
	ReadSystemPrompt      systemprompt.Reader
	ReadConfig            config.ConfigReader
	GetAgentConfiguration agentvendor.ConfigurationPort
	WriteToAgent          systemprompt.AgentWriter
}

func (u *UseCase) Execute(request SyncRequest) error {
	targets, err := u.resolveTargets(request.TargetAgents)
	if err != nil {
		return err
	}

	prompt, err := u.resolveSystemPrompt(request.SystemPromptOverride)
	if err != nil {
		return err
	}

	for _, target := range targets {
		vendorConfig, err := u.GetAgentConfiguration(target)
		if err != nil {
			return ErrUnknownAgent{Name: target}
		}
		if err := u.WriteToAgent(vendorConfig.GlobalInstructionFilePath, prompt); err != nil {
			return err
		}
	}
	return nil
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

func (u *UseCase) resolveSystemPrompt(override *systemprompt.SystemPrompt) (systemprompt.SystemPrompt, error) {
	if override != nil {
		return *override, nil
	}
	return u.ReadSystemPrompt()
}
