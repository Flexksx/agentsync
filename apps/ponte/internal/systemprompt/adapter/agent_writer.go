package adapter

import (
	"os"
	"path/filepath"

	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
)

func WriteToAgent(destinationFilePath string, prompt systemprompt.SystemPrompt) error {
	if err := os.MkdirAll(filepath.Dir(destinationFilePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destinationFilePath, []byte(prompt.Content), 0o644)
}
