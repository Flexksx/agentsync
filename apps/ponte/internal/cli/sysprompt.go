package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	configadapter "github.com/flexksx/ponte/apps/ponte/internal/config/adapter"
	"github.com/flexksx/ponte/apps/ponte/internal/sysprompt"
	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
	promptadapter "github.com/flexksx/ponte/apps/ponte/internal/systemprompt/adapter"
)

func newSyspromptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sysprompt",
		Short: "Show or manage the global system prompt",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := configadapter.ReadConfig()
			if err != nil {
				return err
			}

			prompt, err := promptadapter.ReadSystemPromptFromFile(cfg.SystemPromptFile)
			if errors.Is(err, systemprompt.ErrNoSystemPrompt) {
				cmd.PrintErrln("No system prompt set. Use `ponte sysprompt set <file-or-string>`.")
				return nil
			}
			if err != nil {
				return err
			}

			_, _ = fmt.Fprint(cmd.OutOrStdout(), prompt.Content)
			return nil
		},
	}
	cmd.AddCommand(newSyspromptSetCommand())
	return cmd
}

func newSyspromptSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set <file-or-string>",
		Short: "Persistently set the global system prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := resolveContent(args[0])
			if err != nil {
				return fmt.Errorf("reading system prompt: %w", err)
			}

			cfg, err := configadapter.ReadConfig()
			if err != nil {
				return err
			}

			useCase := &sysprompt.SetUseCase{
				WriteSystemPrompt: func(prompt systemprompt.SystemPrompt) error {
					return promptadapter.WriteSystemPromptToFile(cfg.SystemPromptFile, prompt)
				},
			}

			if err := useCase.Execute(sysprompt.SetRequest{Content: content}); err != nil {
				return err
			}

			cmd.Println("System prompt updated.")
			return nil
		},
	}
}
