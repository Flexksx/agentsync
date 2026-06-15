package sync

import (
	"errors"
	"testing"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	"github.com/flexksx/ponte/apps/ponte/internal/config"
	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
)

func workingUseCase() UseCase {
	return UseCase{
		ReadSystemPrompt: func() (systemprompt.SystemPrompt, error) {
			return systemprompt.SystemPrompt{Content: "default"}, nil
		},
		ReadConfig: func() (config.Config, error) {
			return config.Config{
				Agents: map[agentvendor.AgentVendorName]config.AgentEntry{
					agentvendor.ClaudeCode: {Enabled: true},
				},
			}, nil
		},
		GetAgentConfiguration: func(name agentvendor.AgentVendorName) (agentvendor.AgentVendorConfiguration, error) {
			return agentvendor.AgentVendorConfiguration{
				VendorName:                name,
				GlobalInstructionFilePath: "/fake/" + string(name),
			}, nil
		},
		WriteToAgent: func(_ string, _ systemprompt.SystemPrompt) error {
			return nil
		},
	}
}

func TestExecute_WithExplicitTargets_SkipsConfig(t *testing.T) {
	t.Parallel()
	useCase := workingUseCase()
	configCalled := false
	useCase.ReadConfig = func() (config.Config, error) {
		configCalled = true
		return config.Config{}, nil
	}

	err := useCase.Execute(SyncRequest{
		TargetAgents: []agentvendor.AgentVendorName{agentvendor.ClaudeCode},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if configCalled {
		t.Error("ReadConfig must not be called when targets are explicit")
	}
}

func TestExecute_WithNoTargets_UsesEnabledAgentsFromConfig(t *testing.T) {
	t.Parallel()
	writtenPaths := map[string]bool{}
	useCase := workingUseCase()
	useCase.ReadConfig = func() (config.Config, error) {
		return config.Config{
			Agents: map[agentvendor.AgentVendorName]config.AgentEntry{
				agentvendor.ClaudeCode: {Enabled: true},
				agentvendor.Codex:      {Enabled: false},
				agentvendor.GeminiCLI:  {Enabled: true},
			},
		}, nil
	}
	useCase.WriteToAgent = func(path string, _ systemprompt.SystemPrompt) error {
		writtenPaths[path] = true
		return nil
	}

	err := useCase.Execute(SyncRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(writtenPaths) != 2 {
		t.Errorf("expected 2 writes, got %d: %v", len(writtenPaths), writtenPaths)
	}
	if !writtenPaths["/fake/claude-code"] {
		t.Error("expected write to claude-code")
	}
	if !writtenPaths["/fake/gemini-cli"] {
		t.Error("expected write to gemini-cli")
	}
	if writtenPaths["/fake/codex"] {
		t.Error("must not write to disabled agent codex")
	}
}

func TestExecute_WithNoTargets_NoEnabledAgents_ReturnsErrNoAgentsConfigured(t *testing.T) {
	t.Parallel()
	useCase := workingUseCase()
	useCase.ReadConfig = func() (config.Config, error) {
		return config.Config{
			Agents: map[agentvendor.AgentVendorName]config.AgentEntry{
				agentvendor.ClaudeCode: {Enabled: false},
			},
		}, nil
	}

	err := useCase.Execute(SyncRequest{})

	var target ErrNoAgentsConfigured
	if !errors.As(err, &target) {
		t.Errorf("expected ErrNoAgentsConfigured, got %T: %v", err, err)
	}
}

func TestExecute_WithPromptOverride_WritesOverrideAndSkipsStore(t *testing.T) {
	t.Parallel()
	useCase := workingUseCase()
	storeCalled := false
	useCase.ReadSystemPrompt = func() (systemprompt.SystemPrompt, error) {
		storeCalled = true
		return systemprompt.SystemPrompt{Content: "stored"}, nil
	}
	var writtenPrompt systemprompt.SystemPrompt
	useCase.WriteToAgent = func(_ string, prompt systemprompt.SystemPrompt) error {
		writtenPrompt = prompt
		return nil
	}
	override := systemprompt.SystemPrompt{Content: "override"}

	err := useCase.Execute(SyncRequest{
		TargetAgents:         []agentvendor.AgentVendorName{agentvendor.ClaudeCode},
		SystemPromptOverride: &override,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if storeCalled {
		t.Error("ReadSystemPrompt must not be called when override is provided")
	}
	if writtenPrompt.Content != "override" {
		t.Errorf("expected override content, got %q", writtenPrompt.Content)
	}
}

func TestExecute_WithoutPromptOverride_UsesStoredPrompt(t *testing.T) {
	t.Parallel()
	useCase := workingUseCase()
	useCase.ReadSystemPrompt = func() (systemprompt.SystemPrompt, error) {
		return systemprompt.SystemPrompt{Content: "stored"}, nil
	}
	var writtenPrompt systemprompt.SystemPrompt
	useCase.WriteToAgent = func(_ string, prompt systemprompt.SystemPrompt) error {
		writtenPrompt = prompt
		return nil
	}

	err := useCase.Execute(SyncRequest{
		TargetAgents: []agentvendor.AgentVendorName{agentvendor.ClaudeCode},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if writtenPrompt.Content != "stored" {
		t.Errorf("expected stored content, got %q", writtenPrompt.Content)
	}
}

func TestExecute_WhenAgentConfigurationFails_ReturnsErrUnknownAgent(t *testing.T) {
	t.Parallel()
	useCase := workingUseCase()
	useCase.GetAgentConfiguration = func(_ agentvendor.AgentVendorName) (agentvendor.AgentVendorConfiguration, error) {
		return agentvendor.AgentVendorConfiguration{}, errors.New("not found")
	}

	err := useCase.Execute(SyncRequest{
		TargetAgents: []agentvendor.AgentVendorName{agentvendor.ClaudeCode},
	})

	var target ErrUnknownAgent
	if !errors.As(err, &target) {
		t.Errorf("expected ErrUnknownAgent, got %T: %v", err, err)
	}
	if target.Name != agentvendor.ClaudeCode {
		t.Errorf("expected agent name %q, got %q", agentvendor.ClaudeCode, target.Name)
	}
}

func TestExecute_WhenWriteFails_PropagatesError(t *testing.T) {
	t.Parallel()
	writeErr := errors.New("disk full")
	useCase := workingUseCase()
	useCase.WriteToAgent = func(_ string, _ systemprompt.SystemPrompt) error {
		return writeErr
	}

	err := useCase.Execute(SyncRequest{
		TargetAgents: []agentvendor.AgentVendorName{agentvendor.ClaudeCode},
	})

	if !errors.Is(err, writeErr) {
		t.Errorf("expected write error to be propagated, got %v", err)
	}
}

func TestExecute_WhenConfigReadFails_PropagatesError(t *testing.T) {
	t.Parallel()
	configErr := errors.New("config read failed")
	useCase := workingUseCase()
	useCase.ReadConfig = func() (config.Config, error) {
		return config.Config{}, configErr
	}

	err := useCase.Execute(SyncRequest{})

	if !errors.Is(err, configErr) {
		t.Errorf("expected config error to be propagated, got %v", err)
	}
}

func TestExecute_WhenSystemPromptReadFails_PropagatesError(t *testing.T) {
	t.Parallel()
	promptErr := errors.New("prompt read failed")
	useCase := workingUseCase()
	useCase.ReadSystemPrompt = func() (systemprompt.SystemPrompt, error) {
		return systemprompt.SystemPrompt{}, promptErr
	}

	err := useCase.Execute(SyncRequest{
		TargetAgents: []agentvendor.AgentVendorName{agentvendor.ClaudeCode},
	})

	if !errors.Is(err, promptErr) {
		t.Errorf("expected prompt error to be propagated, got %v", err)
	}
}

func TestExecute_WithMultipleTargets_WritesToEachAgent(t *testing.T) {
	t.Parallel()
	var writtenPaths []string
	useCase := workingUseCase()
	useCase.WriteToAgent = func(path string, _ systemprompt.SystemPrompt) error {
		writtenPaths = append(writtenPaths, path)
		return nil
	}

	err := useCase.Execute(SyncRequest{
		TargetAgents: []agentvendor.AgentVendorName{agentvendor.ClaudeCode, agentvendor.GeminiCLI},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(writtenPaths) != 2 {
		t.Errorf("expected 2 writes, got %d: %v", len(writtenPaths), writtenPaths)
	}
}
