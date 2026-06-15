package systemprompt

type (
	Reader      func() (SystemPrompt, error)
	Writer      func(SystemPrompt) error
	AgentWriter func(destinationFilePath string, prompt SystemPrompt) error
)
