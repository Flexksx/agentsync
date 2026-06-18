package gc

import "github.com/flexksx/ponte/apps/ponte/internal/store"

type Request struct {
	DryRun bool
}

// Result reports the generations a run removed and those kept because a vendor
// still points at them. On a dry run the same Removed set is reported, but no
// generation is actually deleted.
type Result struct {
	Removed []store.Generation
	Kept    []store.Generation
}
