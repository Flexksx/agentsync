// Package gc removes store generations that no vendor currently points to.
package gc

import (
	"fmt"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	"github.com/flexksx/ponte/apps/ponte/internal/store"
)

type UseCase struct {
	KnownVendors          []agentvendor.AgentVendorName
	GetAgentConfiguration agentvendor.ConfigurationPort
	ReadActiveHash        store.ActiveHashReader
	ListGenerations       store.GenerationLister
	RemoveGeneration      store.GenerationRemover
}

// Execute collects the generation hash every known vendor points at, then
// removes every store generation outside that set. All vendors are considered,
// not only enabled ones, because a disabled vendor's symlink still pins a
// generation that must not be collected.
func (u *UseCase) Execute(request Request) (Result, error) {
	activeHashes, err := u.collectActiveHashes()
	if err != nil {
		return Result{}, err
	}

	generations, err := u.ListGenerations()
	if err != nil {
		return Result{}, err
	}

	var result Result
	for _, generation := range generations {
		if activeHashes[generation.Hash] {
			result.Kept = append(result.Kept, generation)
			continue
		}
		if !request.DryRun {
			if err := u.RemoveGeneration(generation); err != nil {
				return Result{}, fmt.Errorf("removing generation %s: %w", generation.Hash, err)
			}
		}
		result.Removed = append(result.Removed, generation)
	}
	return result, nil
}

func (u *UseCase) collectActiveHashes() (map[string]bool, error) {
	active := make(map[string]bool)
	for _, name := range u.KnownVendors {
		vendorConfig, err := u.GetAgentConfiguration(name)
		if err != nil {
			return nil, fmt.Errorf("resolving vendor %q: %w", name, err)
		}
		hash, ok, err := u.ReadActiveHash(vendorConfig.GlobalInstructionFilePath)
		if err != nil {
			return nil, fmt.Errorf("reading active generation for %q: %w", name, err)
		}
		if ok {
			active[hash] = true
		}
	}
	return active, nil
}
