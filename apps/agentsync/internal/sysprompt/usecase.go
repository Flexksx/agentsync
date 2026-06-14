package sysprompt

import "github.com/flexksx/agentsync/apps/agentsync/internal/systemprompt"

type SetUseCase struct {
	WriteSystemPrompt systemprompt.Writer
}

func (u *SetUseCase) Execute(request SetRequest) error {
	return u.WriteSystemPrompt(systemprompt.SystemPrompt{Content: request.Content})
}
