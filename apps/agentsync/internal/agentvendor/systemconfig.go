package agentvendor

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	skillsDir    = "skills"
	subagentsDir = "subagents"

	claudePackageName      = "claude"
	codexPackageName       = "codex"
	geminiPackageName      = "gemini"
	cursorPackageName      = "cursor"
	cursorAgentPackageName = "cursor-agent"

	claudeInstructionFile = "CLAUDE.md"
	codexInstructionFile  = "instructions.md"
	geminiInstructionFile = "GEMINI.md"
	cursorRulesDir        = "rules"
	cursorInstructionFile = "global.mdc"
)

func GetConfiguration(name AgentVendorName) (AgentVendorConfiguration, error) {
	configs, err := platformConfigurations()
	if err != nil {
		return AgentVendorConfiguration{}, err
	}
	cfg, ok := configs[name]
	if !ok {
		return AgentVendorConfiguration{}, &VendorConfigurationNotFoundError{Name: name}
	}
	return cfg, nil
}

func platformConfigurations() (map[AgentVendorName]AgentVendorConfiguration, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		return posixConfigurations(), nil
	case "windows":
		return windowsConfigurations(), nil
	default:
		return nil, &UnsupportedPlatformError{Platform: runtime.GOOS}
	}
}

func posixConfigurations() map[AgentVendorName]AgentVendorConfiguration {
	home, _ := os.UserHomeDir()
	claudeRoot := filepath.Join(home, ".claude")
	codexRoot := filepath.Join(home, ".codex")
	geminiRoot := filepath.Join(home, ".gemini")
	cursorRoot := filepath.Join(home, ".cursor")
	return map[AgentVendorName]AgentVendorConfiguration{
		ClaudeCode: {
			VendorName:                ClaudeCode,
			PackageName:               claudePackageName,
			GlobalInstructionFilePath: filepath.Join(claudeRoot, claudeInstructionFile),
			SkillsDirectoryPath:       filepath.Join(claudeRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(claudeRoot, subagentsDir),
		},
		Codex: {
			VendorName:                Codex,
			PackageName:               codexPackageName,
			GlobalInstructionFilePath: filepath.Join(codexRoot, codexInstructionFile),
			SkillsDirectoryPath:       filepath.Join(codexRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(codexRoot, subagentsDir),
		},
		GeminiCLI: {
			VendorName:                GeminiCLI,
			PackageName:               geminiPackageName,
			GlobalInstructionFilePath: filepath.Join(geminiRoot, geminiInstructionFile),
			SkillsDirectoryPath:       filepath.Join(geminiRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(geminiRoot, subagentsDir),
		},
		CursorAgent: {
			VendorName:                CursorAgent,
			PackageName:               cursorPackageName,
			GlobalInstructionFilePath: filepath.Join(cursorRoot, cursorRulesDir, cursorInstructionFile),
			SkillsDirectoryPath:       filepath.Join(cursorRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(cursorRoot, subagentsDir),
		},
	}
}

func windowsConfigurations() map[AgentVendorName]AgentVendorConfiguration {
	home, _ := os.UserHomeDir()
	roaming := filepath.Join(home, "AppData", "Roaming")
	claudeRoot := filepath.Join(roaming, "Claude")
	codexRoot := filepath.Join(roaming, "Codex")
	geminiRoot := filepath.Join(roaming, "Gemini")
	cursorRoot := filepath.Join(roaming, "Cursor")
	return map[AgentVendorName]AgentVendorConfiguration{
		ClaudeCode: {
			VendorName:                ClaudeCode,
			PackageName:               claudePackageName,
			GlobalInstructionFilePath: filepath.Join(claudeRoot, claudeInstructionFile),
			SkillsDirectoryPath:       filepath.Join(claudeRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(claudeRoot, subagentsDir),
		},
		Codex: {
			VendorName:                Codex,
			PackageName:               codexPackageName,
			GlobalInstructionFilePath: filepath.Join(codexRoot, codexInstructionFile),
			SkillsDirectoryPath:       filepath.Join(codexRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(codexRoot, subagentsDir),
		},
		GeminiCLI: {
			VendorName:                GeminiCLI,
			PackageName:               geminiPackageName,
			GlobalInstructionFilePath: filepath.Join(geminiRoot, geminiInstructionFile),
			SkillsDirectoryPath:       filepath.Join(geminiRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(geminiRoot, subagentsDir),
		},
		CursorAgent: {
			VendorName:                CursorAgent,
			PackageName:               cursorAgentPackageName,
			GlobalInstructionFilePath: filepath.Join(cursorRoot, cursorRulesDir, cursorInstructionFile),
			SkillsDirectoryPath:       filepath.Join(cursorRoot, skillsDir),
			SubagentsDirectoryPath:    filepath.Join(cursorRoot, subagentsDir),
		},
	}
}
