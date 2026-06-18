package adapter

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/flexksx/ponte/apps/ponte/internal/store"
)

const buildDirSuffix = ".build"

// NewLister returns every completed generation in the store. In-progress
// `.build` temp directories and stray files are skipped. A missing store
// directory is not an error — it just means nothing has been synced yet.
func NewLister(storeDir string) store.GenerationLister {
	return func() ([]store.Generation, error) {
		entries, err := os.ReadDir(storeDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil
			}
			return nil, err
		}
		var generations []store.Generation
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasSuffix(entry.Name(), buildDirSuffix) {
				continue
			}
			generations = append(generations, store.Generation{
				Hash:     entry.Name(),
				RootPath: filepath.Join(storeDir, entry.Name()),
			})
		}
		return generations, nil
	}
}

// NewActiveHashReader resolves the generation hash a vendor instruction symlink
// points at. A missing file, a plain (non-symlink) file, or a target outside
// the store all yield ok=false rather than an error: those simply mean the
// vendor is not currently activated by ponte.
func NewActiveHashReader(storeDir string) store.ActiveHashReader {
	return func(instructionFilePath string) (string, bool, error) {
		target, err := os.Readlink(instructionFilePath)
		if err != nil {
			return "", false, nil //nolint:nilerr // absent or non-symlink path means "not activated", not a failure
		}
		relative, err := filepath.Rel(storeDir, target)
		if err != nil {
			return "", false, nil //nolint:nilerr // a target outside the store means "not activated by ponte"
		}
		hash, _, _ := strings.Cut(relative, string(os.PathSeparator))
		if hash == "" || hash == "." || strings.HasPrefix(hash, "..") {
			return "", false, nil
		}
		return hash, true, nil
	}
}

// RemoveGeneration deletes a generation. Generations are made read-only on
// build, so write permissions are restored before removal.
func RemoveGeneration(gen store.Generation) error {
	_ = filepath.WalkDir(gen.RootPath, func(path string, d fs.DirEntry, _ error) error {
		if d != nil && d.IsDir() {
			_ = os.Chmod(path, 0o755)
		} else {
			_ = os.Chmod(path, 0o644)
		}
		return nil
	})
	return os.RemoveAll(gen.RootPath)
}
