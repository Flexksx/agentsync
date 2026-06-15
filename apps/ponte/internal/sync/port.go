package sync

import (
	"fmt"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
)

type ErrNoAgentsConfigured struct{}

func (e ErrNoAgentsConfigured) Error() string {
	return "no agents enabled in config — run with -a to specify agents or edit ~/.config/ponte/config.toml"
}

type ErrUnknownAgent struct {
	Name agentvendor.AgentVendorName
}

func (e ErrUnknownAgent) Error() string {
	return fmt.Sprintf("unknown agent: %s", e.Name)
}
