package cli

import "os"

func resolveContent(argument string) (string, error) {
	data, err := os.ReadFile(argument)
	if err == nil {
		return string(data), nil
	}
	if os.IsNotExist(err) {
		return argument, nil
	}
	return "", err
}
