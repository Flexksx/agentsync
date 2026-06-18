package cli

const shortHashLength = 12

// shortHash truncates a generation hash for display. Full hashes are 32 hex
// chars; the leading 12 are plenty to identify a generation at a glance.
func shortHash(hash string) string {
	if len(hash) <= shortHashLength {
		return hash
	}
	return hash[:shortHashLength]
}
