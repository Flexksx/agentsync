package adapter

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/flexksx/agentsync/apps/agentsync/internal/config"
)

const configFileName = "config.toml"

func ConfigDirectoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "agentsync"), nil
}

func ReadConfig() (config.Config, error) {
	dir, err := ConfigDirectoryPath()
	if err != nil {
		return config.Config{}, err
	}
	path := filepath.Join(dir, configFileName)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return config.Config{}, config.ErrConfigNotInitialized
	}
	var cfg config.Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return config.Config{}, err
	}
	return cfg, nil
}

func WriteConfig(cfg config.Config) error {
	dir, err := ConfigDirectoryPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return err
	}
	path := filepath.Join(dir, configFileName)
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
