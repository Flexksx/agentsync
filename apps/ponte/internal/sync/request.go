package sync

import (
	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
)

type SyncRequest struct {
	SystemPromptOverride *systemprompt.SystemPrompt
	TargetAgents         []agentvendor.AgentVendorName
}
