package adapter

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
)

const (
	skillsDir = "skills"
	agentsDir = "agents"

	claudePackageName      = "claude"
	codexPackageName       = "codex"
	antigravityPackageName = "agy"
	cursorPackageName      = "cursor"
	cursorAgentPackageName = "cursor-agent"
	opencodePackageName    = "opencode"
	piPackageName          = "pi"

	claudeInstructionFile      = "CLAUDE.md"
	codexInstructionFile       = "instructions.md"
	geminiInstructionFile      = "GEMINI.md"
	cursorRulesDir             = "rules"
	cursorInstructionFile      = "global.mdc"
	opencodeInstructionFile    = "AGENTS.md"
	piInstructionFile          = "AGENTS.md"
)

func GetConfiguration(name agentvendor.AgentVendorName) (agentvendor.AgentVendorConfiguration, error) {
	configs, err := platformConfigurations()
	if err != nil {
		return agentvendor.AgentVendorConfiguration{}, err
	}
	cfg, ok := configs[name]
	if !ok {
		return agentvendor.AgentVendorConfiguration{}, &agentvendor.VendorConfigurationNotFoundError{Name: name}
	}
	return cfg, nil
}

func platformConfigurations() (map[agentvendor.AgentVendorName]agentvendor.AgentVendorConfiguration, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		return posixConfigurations(), nil
	case "windows":
		return windowsConfigurations(), nil
	default:
		return nil, &agentvendor.UnsupportedPlatformError{Platform: runtime.GOOS}
	}
}

func posixConfigurations() map[agentvendor.AgentVendorName]agentvendor.AgentVendorConfiguration {
	home, _ := os.UserHomeDir()
	claudeRoot := filepath.Join(home, ".claude")
	codexRoot := filepath.Join(home, ".codex")
	geminiRoot := filepath.Join(home, ".gemini")
	antigravityRoot := filepath.Join(geminiRoot, "antigravity-cli")
	cursorRoot := filepath.Join(home, ".cursor")
	opencodeRoot := filepath.Join(home, ".config", "opencode")
	piRoot := filepath.Join(home, ".pi", "agent")
	return map[agentvendor.AgentVendorName]agentvendor.AgentVendorConfiguration{
		agentvendor.ClaudeCode: {
			VendorName:                agentvendor.ClaudeCode,
			PackageName:               claudePackageName,
			GlobalInstructionFilePath: filepath.Join(claudeRoot, claudeInstructionFile),
			SkillsDirectoryPath:       filepath.Join(claudeRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(claudeRoot, agentsDir),
		},
		agentvendor.Codex: {
			VendorName:                agentvendor.Codex,
			PackageName:               codexPackageName,
			GlobalInstructionFilePath: filepath.Join(codexRoot, codexInstructionFile),
			SkillsDirectoryPath:       filepath.Join(codexRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(codexRoot, agentsDir),
		},
		agentvendor.AntigravityCLI: {
			VendorName:                agentvendor.AntigravityCLI,
			PackageName:               antigravityPackageName,
			GlobalInstructionFilePath: filepath.Join(geminiRoot, geminiInstructionFile),
			SkillsDirectoryPath:       filepath.Join(antigravityRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(antigravityRoot, agentsDir),
		},
		agentvendor.CursorAgent: {
			VendorName:                agentvendor.CursorAgent,
			PackageName:               cursorPackageName,
			GlobalInstructionFilePath: filepath.Join(cursorRoot, cursorRulesDir, cursorInstructionFile),
			SkillsDirectoryPath:       filepath.Join(cursorRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(cursorRoot, agentsDir),
		},
		agentvendor.OpenCode: {
			VendorName:                agentvendor.OpenCode,
			PackageName:               opencodePackageName,
			GlobalInstructionFilePath: filepath.Join(opencodeRoot, opencodeInstructionFile),
			SkillsDirectoryPath:       filepath.Join(opencodeRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(opencodeRoot, agentsDir),
		},
		agentvendor.PiAgent: {
			VendorName:                agentvendor.PiAgent,
			PackageName:               piPackageName,
			GlobalInstructionFilePath: filepath.Join(piRoot, piInstructionFile),
			SkillsDirectoryPath:       filepath.Join(piRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(piRoot, agentsDir),
		},
	}
}

func windowsConfigurations() map[agentvendor.AgentVendorName]agentvendor.AgentVendorConfiguration {
	home, _ := os.UserHomeDir()
	roaming := filepath.Join(home, "AppData", "Roaming")
	claudeRoot := filepath.Join(roaming, "Claude")
	codexRoot := filepath.Join(roaming, "Codex")
	geminiRoot := filepath.Join(roaming, "Gemini")
	antigravityRoot := filepath.Join(geminiRoot, "antigravity-cli")
	cursorRoot := filepath.Join(roaming, "Cursor")
	opencodeRoot := filepath.Join(roaming, "opencode")
	piRoot := filepath.Join(home, ".pi", "agent")
	return map[agentvendor.AgentVendorName]agentvendor.AgentVendorConfiguration{
		agentvendor.ClaudeCode: {
			VendorName:                agentvendor.ClaudeCode,
			PackageName:               claudePackageName,
			GlobalInstructionFilePath: filepath.Join(claudeRoot, claudeInstructionFile),
			SkillsDirectoryPath:       filepath.Join(claudeRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(claudeRoot, agentsDir),
		},
		agentvendor.Codex: {
			VendorName:                agentvendor.Codex,
			PackageName:               codexPackageName,
			GlobalInstructionFilePath: filepath.Join(codexRoot, codexInstructionFile),
			SkillsDirectoryPath:       filepath.Join(codexRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(codexRoot, agentsDir),
		},
		agentvendor.AntigravityCLI: {
			VendorName:                agentvendor.AntigravityCLI,
			PackageName:               antigravityPackageName,
			GlobalInstructionFilePath: filepath.Join(geminiRoot, geminiInstructionFile),
			SkillsDirectoryPath:       filepath.Join(antigravityRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(antigravityRoot, agentsDir),
		},
		agentvendor.CursorAgent: {
			VendorName:                agentvendor.CursorAgent,
			PackageName:               cursorAgentPackageName,
			GlobalInstructionFilePath: filepath.Join(cursorRoot, cursorRulesDir, cursorInstructionFile),
			SkillsDirectoryPath:       filepath.Join(cursorRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(cursorRoot, agentsDir),
		},
		agentvendor.OpenCode: {
			VendorName:                agentvendor.OpenCode,
			PackageName:               opencodePackageName,
			GlobalInstructionFilePath: filepath.Join(opencodeRoot, opencodeInstructionFile),
			SkillsDirectoryPath:       filepath.Join(opencodeRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(opencodeRoot, agentsDir),
		},
		agentvendor.PiAgent: {
			VendorName:                agentvendor.PiAgent,
			PackageName:               piPackageName,
			GlobalInstructionFilePath: filepath.Join(piRoot, piInstructionFile),
			SkillsDirectoryPath:       filepath.Join(piRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(piRoot, agentsDir),
		},
	}
}
