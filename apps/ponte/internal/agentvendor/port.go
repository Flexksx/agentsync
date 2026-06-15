package agentvendor

import "fmt"

type ConfigurationPort func(AgentVendorName) (AgentVendorConfiguration, error)

type VendorConfigurationNotFoundError struct {
	Name AgentVendorName
}

func (e *VendorConfigurationNotFoundError) Error() string {
	return fmt.Sprintf("no configuration found for vendor: %s", e.Name)
}

type UnsupportedPlatformError struct {
	Platform string
}

func (e *UnsupportedPlatformError) Error() string {
	return fmt.Sprintf("unsupported platform: %s", e.Platform)
}
