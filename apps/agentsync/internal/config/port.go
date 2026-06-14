package config

import "errors"

type (
	ConfigReader func() (Config, error)
	ConfigWriter func(Config) error
)

var ErrConfigNotInitialized = errors.New("agentsync config not initialized")
