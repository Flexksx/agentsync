package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	vendoradapter "github.com/flexksx/ponte/apps/ponte/internal/agentvendor/adapter"
	"github.com/flexksx/ponte/apps/ponte/internal/gc"
	storeadapter "github.com/flexksx/ponte/apps/ponte/internal/store/adapter"
)

func newGcCommand() *cobra.Command {
	var dryRunFlag bool

	cmd := &cobra.Command{
		Use:   "gc",
		Short: "Remove store generations no vendor points to",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			storeDir, err := storeadapter.StoreDirectoryPath()
			if err != nil {
				return fmt.Errorf("resolving store directory: %w", err)
			}

			useCase := &gc.UseCase{
				KnownVendors:          agentvendor.AllVendorNames(),
				GetAgentConfiguration: vendoradapter.GetConfiguration,
				ReadActiveHash:        storeadapter.NewActiveHashReader(storeDir),
				ListGenerations:       storeadapter.NewLister(storeDir),
				RemoveGeneration:      storeadapter.RemoveGeneration,
			}

			result, err := useCase.Execute(gc.Request{DryRun: dryRunFlag})
			if err != nil {
				return err
			}

			printGcResult(cmd, result, dryRunFlag)
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show which generations would be removed without deleting them")

	return cmd
}

func printGcResult(cmd *cobra.Command, result gc.Result, dryRun bool) {
	out := cmd.OutOrStdout()
	if len(result.Removed) == 0 {
		_, _ = fmt.Fprintf(out, "Nothing to remove; %d generation(s) in use.\n", len(result.Kept))
		return
	}

	verb := "Removed"
	if dryRun {
		verb = "Would remove"
	}
	_, _ = fmt.Fprintf(out, "%s %d generation(s), kept %d in use:\n", verb, len(result.Removed), len(result.Kept))
	for _, generation := range result.Removed {
		_, _ = fmt.Fprintf(out, "  %s\n", shortHash(generation.Hash))
	}
}
