package sync

import (
	"github.com/flexksx/agentsync/apps/agentsync/internal/agentvendor"
	"github.com/flexksx/agentsync/apps/agentsync/internal/systemprompt"
)

type SyncRequest struct {
	SystemPromptOverride *systemprompt.SystemPrompt
	TargetAgents         []agentvendor.AgentVendorName
}
