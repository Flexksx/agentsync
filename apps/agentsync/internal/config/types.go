package config

import "github.com/flexksx/agentsync/apps/agentsync/internal/agentvendor"

type AgentEntry struct {
	Enabled bool `toml:"enabled"`
}

type Config struct {
	Agents map[agentvendor.AgentVendorName]AgentEntry `toml:"agents"`
}

func DefaultConfig() Config {
	return Config{
		Agents: map[agentvendor.AgentVendorName]AgentEntry{
			agentvendor.ClaudeCode:  {Enabled: true},
			agentvendor.Codex:       {Enabled: true},
			agentvendor.GeminiCLI:   {Enabled: true},
			agentvendor.CursorAgent: {Enabled: true},
		},
	}
}
