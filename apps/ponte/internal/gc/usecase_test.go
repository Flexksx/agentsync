package gc

import (
	"errors"
	"testing"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	"github.com/flexksx/ponte/apps/ponte/internal/store"
)

// activeFor maps a vendor's instruction path to the generation hash its symlink
// resolves to, so tests declare "vendor X points at generation Y" directly.
func newUseCase(active map[string]string, generations []store.Generation, removed *[]string) UseCase {
	return UseCase{
		KnownVendors: agentvendor.AllVendorNames(),
		GetAgentConfiguration: func(name agentvendor.AgentVendorName) (agentvendor.AgentVendorConfiguration, error) {
			return agentvendor.AgentVendorConfiguration{
				VendorName:                name,
				GlobalInstructionFilePath: "/fake/" + string(name) + "/instruction",
			}, nil
		},
		ReadActiveHash: func(instructionFilePath string) (string, bool, error) {
			hash, ok := active[instructionFilePath]
			return hash, ok, nil
		},
		ListGenerations: func() ([]store.Generation, error) {
			return generations, nil
		},
		RemoveGeneration: func(gen store.Generation) error {
			*removed = append(*removed, gen.Hash)
			return nil
		},
	}
}

func instructionPath(name agentvendor.AgentVendorName) string {
	return "/fake/" + string(name) + "/instruction"
}

func TestExecute_RemovesGenerationsNoVendorPointsTo(t *testing.T) {
	t.Parallel()
	active := map[string]string{
		instructionPath(agentvendor.ClaudeCode): "keep",
	}
	generations := []store.Generation{{Hash: "keep"}, {Hash: "stale1"}, {Hash: "stale2"}}
	var removed []string

	useCase := newUseCase(active, generations, &removed)
	result, err := useCase.Execute(Request{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 2 {
		t.Errorf("expected 2 removals, got %v", removed)
	}
	if len(result.Kept) != 1 || result.Kept[0].Hash != "keep" {
		t.Errorf("expected to keep 'keep', got %v", result.Kept)
	}
}

func TestExecute_KeepsGenerationsPinnedByDisabledVendors(t *testing.T) {
	t.Parallel()
	// Two different vendors point at two different generations. Even though only
	// one would be enabled in config, GC considers all vendors, so both survive.
	active := map[string]string{
		instructionPath(agentvendor.ClaudeCode): "gen-a",
		instructionPath(agentvendor.Codex):      "gen-b",
	}
	generations := []store.Generation{{Hash: "gen-a"}, {Hash: "gen-b"}, {Hash: "orphan"}}
	var removed []string

	useCase := newUseCase(active, generations, &removed)
	result, err := useCase.Execute(Request{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 1 || removed[0] != "orphan" {
		t.Errorf("expected only 'orphan' removed, got %v", removed)
	}
	if len(result.Kept) != 2 {
		t.Errorf("expected 2 kept, got %v", result.Kept)
	}
}

func TestExecute_DryRun_ReportsButDoesNotRemove(t *testing.T) {
	t.Parallel()
	active := map[string]string{instructionPath(agentvendor.ClaudeCode): "keep"}
	generations := []store.Generation{{Hash: "keep"}, {Hash: "stale"}}
	var removed []string

	useCase := newUseCase(active, generations, &removed)
	result, err := useCase.Execute(Request{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 0 {
		t.Errorf("dry run must not remove anything, removed %v", removed)
	}
	if len(result.Removed) != 1 || result.Removed[0].Hash != "stale" {
		t.Errorf("expected 'stale' reported as removable, got %v", result.Removed)
	}
}

func TestExecute_WhenRemovalFails_PropagatesError(t *testing.T) {
	t.Parallel()
	removeErr := errors.New("permission denied")
	useCase := newUseCase(map[string]string{}, []store.Generation{{Hash: "stale"}}, &[]string{})
	useCase.RemoveGeneration = func(_ store.Generation) error {
		return removeErr
	}

	_, err := useCase.Execute(Request{})
	if !errors.Is(err, removeErr) {
		t.Errorf("expected removal error to be propagated, got %v", err)
	}
}

func TestExecute_WhenListFails_PropagatesError(t *testing.T) {
	t.Parallel()
	listErr := errors.New("store unreadable")
	useCase := newUseCase(map[string]string{}, nil, &[]string{})
	useCase.ListGenerations = func() ([]store.Generation, error) {
		return nil, listErr
	}

	_, err := useCase.Execute(Request{})
	if !errors.Is(err, listErr) {
		t.Errorf("expected list error to be propagated, got %v", err)
	}
}
