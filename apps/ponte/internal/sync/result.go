package sync

import "github.com/flexksx/ponte/apps/ponte/internal/agentvendor"

// SyncResult reports what a sync produced so the presentation layer does not
// have to re-derive it. On a dry run GenerationHash is the generation that
// would be built and no vendor is touched.
type SyncResult struct {
	GenerationHash string
	Targets        []agentvendor.AgentVendorName
	DryRun         bool
}
