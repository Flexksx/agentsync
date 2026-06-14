package agentvendor

import "os/exec"

func IsInstalled(name AgentVendorName, getConfiguration ConfigurationPort) (bool, error) {
	configuration, err := getConfiguration(name)
	if err != nil {
		return false, err
	}
	_, err = exec.LookPath(configuration.PackageName)
	return err == nil, nil
}
