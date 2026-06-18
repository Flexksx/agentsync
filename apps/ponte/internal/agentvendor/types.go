// Package agentvendor defines supported AI agent vendors and their configuration.
package agentvendor

type AgentVendorName string

const (
	ClaudeCode  AgentVendorName = "claude-code"
	Codex       AgentVendorName = "codex"
	GeminiCLI   AgentVendorName = "gemini-cli"
	CursorAgent AgentVendorName = "cursor-agent"
)

// AllVendorNames lists every supported vendor. Consumers that must consider
// vendors regardless of config (garbage collection, status) iterate this so a
// disabled vendor still holding store symlinks is not overlooked.
func AllVendorNames() []AgentVendorName {
	return []AgentVendorName{ClaudeCode, Codex, GeminiCLI, CursorAgent}
}

type AgentVendorConfiguration struct {
	VendorName                AgentVendorName
	PackageName               string
	GlobalInstructionFilePath string
	SkillsDirectoryPath       string
	SubagentsDirectoryPath    string
}
