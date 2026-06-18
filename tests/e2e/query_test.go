package e2e

import (
	"strings"
	"testing"
)

// `ponte skills` lists each declared skill with its name and resolved source.
func TestSkills_ListsDeclaredSkills(t *testing.T) {
	h := newHarness(t)
	h.bootstrap()

	skillFixtureDir := repoFixtureDir(t, "simple_skill")
	writeConfigWithLocalSkill(t, h, "simple-skill", skillFixtureDir)

	stdout, _ := h.mustRun("skills")

	for _, want := range []string{"NAME", "simple-skill", "local", skillFixtureDir} {
		if !strings.Contains(stdout, want) {
			t.Errorf("expected `skills` output to contain %q, got:\n%s", want, stdout)
		}
	}
}

// `ponte skills` with nothing declared reports an empty state rather than a
// blank or error.
func TestSkills_NoneConfigured_ReportsEmpty(t *testing.T) {
	h := newHarness(t)
	h.bootstrap()

	stdout, _ := h.mustRun("skills")

	if !strings.Contains(stdout, "No skills configured.") {
		t.Errorf("expected empty-state message, got:\n%s", stdout)
	}
}

// Bare `ponte sysprompt` prints the stored prompt verbatim to stdout so it can
// be piped or redirected.
func TestSysprompt_Show_PrintsStoredPrompt(t *testing.T) {
	h := newHarness(t)
	h.bootstrap()
	h.mustRun("sysprompt", "set", samplePrompt)

	stdout, _ := h.mustRun("sysprompt")

	if stdout != samplePrompt {
		t.Errorf("expected sysprompt show to print stored prompt\n--- want ---\n%q\n--- got ---\n%q", samplePrompt, stdout)
	}
}

// `ponte sysprompt set` still mutates while bare `sysprompt` shows — the two
// must round-trip.
func TestSysprompt_SetThenShow_RoundTrips(t *testing.T) {
	h := newHarness(t)
	h.bootstrap()

	updated := "# Updated prompt\n\nNew instructions.\n"
	h.mustRun("sysprompt", "set", updated)

	stdout, _ := h.mustRun("sysprompt")
	if stdout != updated {
		t.Errorf("expected updated prompt, got:\n%s", stdout)
	}
}
