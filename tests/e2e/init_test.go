package e2e

import (
	"os"
	"strings"
	"testing"
)

// Running sync against an empty $HOME should bootstrap a default config.toml
// and an empty system prompt file, then proceed with the sync. The behaviour is
// documented in cli/sync.go's ensureConfigInitialized.
func TestInit_FirstSyncBootstrapsConfig(t *testing.T) {
	h := newHarness(t)

	stdout, stderr := h.mustRun("sync")

	// cobra's cmd.Print* writes to stderr by default; either stream is fine.
	combined := stdout + stderr
	if !strings.Contains(combined, "Initialized agentsync config") {
		t.Errorf("expected init banner, got stdout=%q stderr=%q", stdout, stderr)
	}

	// config.toml must exist with all four vendors enabled.
	cfg := h.readFile(h.configFile())
	for _, name := range []string{"claude-code", "codex", "gemini-cli", "cursor-agent"} {
		if !strings.Contains(cfg, name) {
			t.Errorf("config.toml missing vendor %q:\n%s", name, cfg)
		}
	}

	// Default system prompt file is AGENTS.md and starts empty.
	if got := h.readFile(h.promptFile("AGENTS.md")); got != "" {
		t.Errorf("expected empty AGENTS.md after init, got %q", got)
	}

	// Empty prompt content is propagated to every vendor instruction file.
	for vendor, path := range h.vendorPaths() {
		h.assertFileEquals(path, "")
		_ = vendor
	}
}

// A second sync after init must not error or duplicate the init banner — the
// bootstrap path is a one-shot.
func TestInit_IsIdempotent(t *testing.T) {
	h := newHarness(t)

	h.mustRun("sync")
	stdout, stderr := h.mustRun("sync")

	combined := stdout + stderr
	if strings.Contains(combined, "Initialized agentsync config") {
		t.Errorf("second sync should not re-run init, stdout=%q stderr=%q", stdout, stderr)
	}
}

// If $HOME points at an unwritable location, sync must surface the error
// instead of silently no-oping. We don't try to simulate that on Windows where
// chmod semantics differ; the assertion is Unix-only.
func TestInit_FailsWhenHomeUnwritable(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses unix permissions")
	}
	if isWindows() {
		t.Skip("permission semantics differ on Windows")
	}

	h := newHarness(t)
	if err := os.Chmod(h.home, 0o500); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(h.home, 0o700) })

	_, _, err := h.run("sync")
	if err == nil {
		t.Fatal("expected sync to fail with unwritable $HOME")
	}
}
