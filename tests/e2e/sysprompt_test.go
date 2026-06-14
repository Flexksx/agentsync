package e2e

import (
	"strings"
	"testing"
)

// `sysprompt set "literal string"` writes the literal to the configured prompt
// file. The argument is treated as an inline string when no such file exists.
func TestSyspromptSet_InlineString(t *testing.T) {
	h := newHarness(t)
	h.mustRun("sync") // bootstrap

	const content = "literal inline content with spaces"
	h.mustRun("sysprompt", "set", content)

	h.assertFileEquals(h.promptFile("AGENTS.md"), content)
}

// `sysprompt set <file>` reads the argument as a file when it exists on disk.
// We use a fixture under tests/e2e/fixtures/ so the assertion is reproducible.
func TestSyspromptSet_FileArgument(t *testing.T) {
	h := newHarness(t)
	h.mustRun("sync") // bootstrap

	fixture := repoFixturePath(t, "simple_prompt.md")
	want := h.readFile(fixture)

	h.mustRun("sysprompt", "set", fixture)
	h.assertFileEquals(h.promptFile("AGENTS.md"), want)
}

// Running set twice in a row overwrites — no appending, no merging.
func TestSyspromptSet_Overwrites(t *testing.T) {
	h := newHarness(t)
	h.mustRun("sync")

	h.mustRun("sysprompt", "set", "first")
	h.mustRun("sysprompt", "set", "second")

	h.assertFileEquals(h.promptFile("AGENTS.md"), "second")
}

// Multi-line content via file argument must round-trip byte-for-byte. We've
// seen line-ending normalization bugs in CLIs that round through string→file→string.
func TestSyspromptSet_MultiLineFile(t *testing.T) {
	h := newHarness(t)
	h.mustRun("sync")

	fixture := repoFixturePath(t, "unicode_prompt.md")
	want := h.readFile(fixture)

	h.mustRun("sysprompt", "set", fixture)
	h.assertFileEquals(h.promptFile("AGENTS.md"), want)

	// And the bytes are preserved over a full sync round-trip.
	h.mustRun("sync")
	for _, p := range h.vendorPaths() {
		h.assertFileEquals(p, want)
	}
}

// `sysprompt set` with zero args is a usage error from cobra (ExactArgs(1)).
func TestSyspromptSet_RequiresArgument(t *testing.T) {
	h := newHarness(t)
	h.mustRun("sync")

	_, stderr, err := h.run("sysprompt", "set")
	if err == nil {
		t.Fatal("expected non-zero exit when arg is missing")
	}
	if !strings.Contains(strings.ToLower(stderr), "arg") {
		t.Errorf("expected usage error mentioning args, got:\n%s", stderr)
	}
}
