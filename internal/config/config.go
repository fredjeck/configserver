// Package config groups all configuration related utilities
package config

import (
	"log/slog"
	"os"
	"strings"
)

// InitLogging sets up logging based on the CONFIGSERVER_ENV environment variable.
// JSON based logging will be automatically enabled if the environment name does not contain the "dev" substring
func InitLogging() {
	env := strings.ToLower(os.Getenv(EnvConfigServerEnvironment))

	var logger *slog.Logger
	if strings.Contains(env, "dev") {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)
}

// LoadFrom loads the configuration file from the provided path
// This method always return a Configuration object with at least the Environment configuration loaded
func LoadFrom(path string) (*Configuration, error) {
	config := &Configuration{
		Environment: &Environment{
			Kind: strings.ToLower(os.Getenv(EnvConfigServerEnvironment)),
		},
	}
	return config, nil
}

// Configuration groups all the supported configuration options
type Configuration struct {
	*Environment
}

// Environment gathers all the environment variable used by ConfigServer
type Environment struct {
	Kind string
}
