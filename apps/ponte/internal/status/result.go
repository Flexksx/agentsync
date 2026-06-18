package status

import "github.com/flexksx/ponte/apps/ponte/internal/agentvendor"

// VendorStatus is one vendor's activation state. ActiveHash is the generation
// its instruction symlink resolves to; HasActive is false when the vendor has
// never been synced (no symlink into the store).
type VendorStatus struct {
	Name       agentvendor.AgentVendorName
	Enabled    bool
	ActiveHash string
	HasActive  bool
}

// Report is the whole-system view: the generation a sync would build now, and
// every vendor's current state to compare against it.
type Report struct {
	WouldBeHash string
	Vendors     []VendorStatus
}

// InSync reports whether the vendor's activated generation already matches what
// a sync would build.
func (v VendorStatus) InSync(wouldBeHash string) bool {
	return v.HasActive && v.ActiveHash == wouldBeHash
}
