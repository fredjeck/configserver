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

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var logger *slog.Logger
	if strings.Contains(env, "dev") {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	slog.SetDefault(logger)
}

// LoadFrom loads the configuration file from the provided path
// This method always return a Configuration object with at least the Environment configuration loaded
func LoadFrom(path string) (*Configuration, error) {

	kind := strings.ToLower(os.Getenv(EnvConfigServerEnvironment))
	if len(kind) == 0 {
		kind = "production"
	}

	keys := strings.ToLower(os.Getenv(EnvConfigServerKeysPath))
	if len(keys) == 0 {
		keys = "/var/run/configserver/keys"
	}

	config := &Configuration{
		Environment: &Environment{
			Kind:     kind,
			KeysPath: keys,
		},
		Server: &Server{
			ListenOn: ":8080",
		},
	}
	return config, nil
}

// Configuration groups all the supported configuration options
type Configuration struct {
	// Environment variables for ease of use
	*Environment
	// Server related settings
	*Server
}

// Environment gathers all the environment variable used by ConfigServer
type Environment struct {
	Kind     string // environment kind i,e dev, int, prod
	KeysPath string // path where the various keys can be found
}

// Server groups all the configserver related settings
type Server struct {
	ListenOn string // address and port on which the server will listen for incoming requests
}

// LogEnvironment logs the current environment configuration
func (c *Configuration) LogEnvironment() {
	slog.Info("Configserver Runtime Environment",
		EnvConfigServerEnvironment, c.Environment.Kind,
		EnvConfigServerKeysPath, c.Environment.KeysPath,
	)
}
