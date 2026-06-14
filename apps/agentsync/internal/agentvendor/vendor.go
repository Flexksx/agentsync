package agentvendor

type AgentVendorName string

const (
	ClaudeCode  AgentVendorName = "claude-code"
	Codex       AgentVendorName = "codex"
	GeminiCLI   AgentVendorName = "gemini-cli"
	CursorAgent AgentVendorName = "cursor-agent"
)

type AgentVendorConfiguration struct {
	VendorName                AgentVendorName
	PackageName               string
	GlobalInstructionFilePath string
	SkillsDirectoryPath       string
	SubagentsDirectoryPath    string
}
