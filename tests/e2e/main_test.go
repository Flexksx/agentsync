// Package e2e contains end-to-end tests that exercise the ponte CLI as a
// real subprocess against an isolated $HOME, so the same suite runs identically
// on Linux, macOS, and Windows CI runners.
package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// binaryPath holds the absolute path to the compiled ponte binary used by
// every test in this package. It is populated by TestMain.
var binaryPath string

// TestMain builds the CLI once into a temporary directory and reuses it across
// all tests in this package. Building once is dramatically faster than letting
// each test shell out to `go run`.
func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "e2e setup failed:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func runMain(m *testing.M) (int, error) {
	tmpDir, err := os.MkdirTemp("", "ponte-e2e-bin-*")
	if err != nil {
		return 0, fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	binName := "ponte"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binaryPath = filepath.Join(tmpDir, binName)

	repoRoot, err := findRepoRoot()
	if err != nil {
		return 0, err
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, "./apps/ponte")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return 0, fmt.Errorf("go build failed: %w\n%s", err, out)
	}

	return m.Run(), nil
}

// findRepoRoot walks up from the current working directory looking for go.mod,
// since `go test` runs each package from its own directory.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}
