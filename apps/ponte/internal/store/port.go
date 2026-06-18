package store

type GenerationBuilder func(input BuildInput) (Generation, error)

// HashComputer derives a generation's content hash from its inputs without
// writing anything to the store. Used to preview a sync and to detect drift.
type HashComputer func(input BuildInput) (string, error)

// VendorActivator symlinks the generation's instruction file, skills, and
// subagents into the vendor-specific directory layout.
type VendorActivator func(gen Generation, instructionFilePath, skillsDirPath, subagentsDirPath string) error

// GenerationLister returns every generation currently present in the store.
type GenerationLister func() ([]Generation, error)

// ActiveHashReader reads the generation hash a vendor's instruction symlink
// resolves to. ok is false when the path is absent or not a store symlink.
type ActiveHashReader func(instructionFilePath string) (hash string, ok bool, err error)

// GenerationRemover deletes a single generation from the store, restoring write
// permissions first since generations are made read-only on build.
type GenerationRemover func(gen Generation) error
