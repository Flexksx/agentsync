package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flexksx/ponte/apps/ponte/internal/agentvendor"
	vendoradapter "github.com/flexksx/ponte/apps/ponte/internal/agentvendor/adapter"
	"github.com/flexksx/ponte/apps/ponte/internal/config"
	configadapter "github.com/flexksx/ponte/apps/ponte/internal/config/adapter"
	skilladapter "github.com/flexksx/ponte/apps/ponte/internal/skill/adapter"
	"github.com/flexksx/ponte/apps/ponte/internal/status"
	storeadapter "github.com/flexksx/ponte/apps/ponte/internal/store/adapter"
	"github.com/flexksx/ponte/apps/ponte/internal/systemprompt"
	promptadapter "github.com/flexksx/ponte/apps/ponte/internal/systemprompt/adapter"
)

const (
	statusVendorHeader  = "VENDOR"
	statusEnabledHeader = "ENABLED"
	statusActiveHeader  = "ACTIVE"
	statusStateHeader   = "STATE"
	statusVendorWidth   = 12
	statusEnabledWidth  = 7
	statusActiveWidth   = 12

	stateInSync    = "in sync"
	stateDrifted   = "drifted"
	stateNotSynced = "not synced"
	stateDisabled  = "disabled"

	noActiveGeneration = "—"
	yesLabel           = "yes"
	noLabel            = "no"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the active generation per vendor and whether sources have drifted",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := configadapter.ReadConfig()
			if err != nil {
				return err
			}

			storeDir, err := storeadapter.StoreDirectoryPath()
			if err != nil {
				return fmt.Errorf("resolving store directory: %w", err)
			}

			gitCacheDir, err := skilladapter.GitCacheDirectoryPath()
			if err != nil {
				return fmt.Errorf("resolving skill cache directory: %w", err)
			}

			useCase := &status.UseCase{
				KnownVendors: agentvendor.AllVendorNames(),
				ReadConfig:   func() (config.Config, error) { return cfg, nil },
				ReadSystemPrompt: func() (systemprompt.SystemPrompt, error) {
					return promptadapter.ReadSystemPromptFromFile(cfg.SystemPromptFile)
				},
				GetAgentConfiguration: vendoradapter.GetConfiguration,
				ResolveSkill:          skilladapter.NewResolver(gitCacheDir),
				ComputeHash:           storeadapter.ComputeHash,
				ReadActiveHash:        storeadapter.NewActiveHashReader(storeDir),
			}

			report, err := useCase.Execute()
			if err != nil {
				return err
			}

			printStatus(cmd, report)
			return nil
		},
	}
}

func printStatus(cmd *cobra.Command, report status.Report) {
	out := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(out, "Would-be generation: %s\n\n", shortHash(report.WouldBeHash))

	_, _ = fmt.Fprintf(out, "%-*s  %-*s  %-*s  %s\n",
		statusVendorWidth, statusVendorHeader,
		statusEnabledWidth, statusEnabledHeader,
		statusActiveWidth, statusActiveHeader,
		statusStateHeader)

	for _, vendor := range report.Vendors {
		_, _ = fmt.Fprintf(out, "%-*s  %-*s  %-*s  %s\n",
			statusVendorWidth, vendor.Name,
			statusEnabledWidth, enabledLabel(vendor.Enabled),
			statusActiveWidth, activeLabel(vendor),
			vendorState(vendor, report.WouldBeHash))
	}
}

func enabledLabel(enabled bool) string {
	if enabled {
		return yesLabel
	}
	return noLabel
}

func activeLabel(vendor status.VendorStatus) string {
	if !vendor.HasActive {
		return noActiveGeneration
	}
	return shortHash(vendor.ActiveHash)
}

// vendorState classifies a vendor for the STATE column. A disabled vendor is
// reported as such regardless of drift because a sync will not touch it.
func vendorState(vendor status.VendorStatus, wouldBeHash string) string {
	if !vendor.Enabled {
		return stateDisabled
	}
	if !vendor.HasActive {
		return stateNotSynced
	}
	if vendor.InSync(wouldBeHash) {
		return stateInSync
	}
	return stateDrifted
}
