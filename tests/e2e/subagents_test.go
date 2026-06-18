package e2e

import (
	"strings"
	"testing"
)

// With no subagents declared, the command must say so rather than printing an
// empty table.
func TestSubagents_EmptyConfig(t *testing.T) {
	h := newHarness(t)
	h.bootstrap()

	stdout, _ := h.mustRun("subagents")
	if !strings.Contains(stdout, "No subagents configured.") {
		t.Errorf("expected empty-state message, got:\n%s", stdout)
	}
}

// A declared subagent must be listed with its name and source.
func TestSubagents_ListsDeclared(t *testing.T) {
	h := newHarness(t)
	h.bootstrap()

	subagentDir := repoFixtureDir(t, "subagents")
	writeConfigWithLocalSubagent(t, h, "my-agents", subagentDir)

	stdout, _ := h.mustRun("subagents")
	if !strings.Contains(stdout, "my-agents") {
		t.Errorf("expected subagent name in output, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "local") {
		t.Errorf("expected source type in output, got:\n%s", stdout)
	}
}
