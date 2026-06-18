package e2e

import (
	"os"
	"strings"
	"testing"
)

// countStoreGenerations returns the number of completed generation directories
// in the store (ignoring in-progress .build temp dirs).
func countStoreGenerations(t *testing.T, h *harness) int {
	t.Helper()
	entries, err := os.ReadDir(h.storePath())
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("read store dir: %v", err)
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasSuffix(entry.Name(), ".build") {
			count++
		}
	}
	return count
}

// After several distinct syncs, only the active generation should survive gc;
// the superseded ones are orphaned and removed.
func TestGc_RemovesOrphanedGenerations(t *testing.T) {
	if isWindows() {
		t.Skip("symlink tests require Unix")
	}
	h := newHarness(t)
	h.bootstrap()

	h.mustRun("sysprompt", "set", "v1")
	h.mustRun("sync")
	h.mustRun("sysprompt", "set", "v2")
	h.mustRun("sync")
	h.mustRun("sysprompt", "set", "v3")
	h.mustRun("sync")

	if got := countStoreGenerations(t, h); got < 2 {
		t.Fatalf("expected multiple generations before gc, got %d", got)
	}

	stdout, _ := h.mustRun("gc")
	if !strings.Contains(stdout, "Removed") {
		t.Errorf("expected gc to report removals, got:\n%s", stdout)
	}

	if got := countStoreGenerations(t, h); got != 1 {
		t.Errorf("expected exactly 1 generation after gc, got %d", got)
	}

	// The active generation must still resolve — vendor files intact.
	h.assertFileEquals(h.vendorPaths()["claude-code"], "v3")
}

// gc --dry-run reports what it would remove but deletes nothing.
func TestGc_DryRun_RemovesNothing(t *testing.T) {
	if isWindows() {
		t.Skip("symlink tests require Unix")
	}
	h := newHarness(t)
	h.bootstrap()

	h.mustRun("sysprompt", "set", "v1")
	h.mustRun("sync")
	h.mustRun("sysprompt", "set", "v2")
	h.mustRun("sync")

	before := countStoreGenerations(t, h)

	stdout, _ := h.mustRun("gc", "--dry-run")
	if !strings.Contains(stdout, "Would remove") {
		t.Errorf("expected dry-run to report 'Would remove', got:\n%s", stdout)
	}

	if after := countStoreGenerations(t, h); after != before {
		t.Errorf("dry run changed generation count: before=%d after=%d", before, after)
	}
}

// A second gc with no intervening sync has nothing left to collect: the first
// gc already pruned everything but the active generation.
func TestGc_NothingToRemove(t *testing.T) {
	if isWindows() {
		t.Skip("symlink tests require Unix")
	}
	h := newHarness(t)
	h.bootstrap()
	h.mustRun("sysprompt", "set", "v1")
	h.mustRun("sync")
	h.mustRun("gc")

	stdout, _ := h.mustRun("gc")
	if !strings.Contains(stdout, "Nothing to remove") {
		t.Errorf("expected 'Nothing to remove' on second gc, got:\n%s", stdout)
	}
}
