// Package config groups all configuration related utilities
package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// EnvConfigServerEnvironment holds the name of the environment variable storing the current active environment type (dev, int, prod ...)
const EnvConfigServerEnvironment string = "CONFIGSERVER_ENV"

// EnvConfigServerHome holds the path to the directory containing the application's config
const EnvConfigServerHome string = "CONFIGSERVER_HOME"

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

	home := strings.ToLower(os.Getenv(EnvConfigServerHome))
	if len(home) == 0 {
		home = "/var/run/configserver"
	}

	configPath := path
	if _, err := os.Stat(path); err != nil {
		configPath = filepath.Join(home, "configserver.yml")
		if _, err := os.Stat(configPath); err != nil {
			return nil, fmt.Errorf("'%s' configserver configuration cannot be found or is not accessible", path)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("'%s' configserver configuration cannot be found or is not accessible : %w", path, err)
	}

	config := &Configuration{
		Environment: &Environment{
			Kind: kind,
			Home: home,
		},
		Server: &Server{
			ListenOn: ":8080",
		},
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, fmt.Errorf("'%s' cannot unmarshal yaml file : %w", path, err)
	}

	if len(config.CertsLocation) == 0 {
		config.CertsLocation = filepath.Join(home, "certs")
	}

	return config, nil
}

// Configuration groups all the supported configuration options
type Configuration struct {
	CertsLocation string `yaml:"certs_location"` // location of keys and certificates
	// Environment variables for ease of use
	*Environment
	// Server related settings
	*Server
	// Git related settings
	*GitConfiguration `yaml:"git"`
}

// Environment gathers all the environment variable used by ConfigServer
type Environment struct {
	Kind string // environment kind i,e dev, int, prod
	Home string // path to the directory containing configserver configuration files
}

// Server groups all the configserver related settings
type Server struct {
	ListenOn string `yaml:"listen_on"` // address and port on which the server will listen for incoming requests
}

type GitConfiguration struct {
	RepositoriesCheckoutLocation      string `yaml:"repositories_checkout_location"`      // Git repositories checkout location
	RepositoriesConfigurationLocation string `yaml:"repositories_configuration_location"` // Path where git configuration files can be found
}

// LogEnvironment logs the current environment configuration
func (c *Configuration) LogEnvironment() {
	slog.Info("Configserver Runtime Environment",
		EnvConfigServerEnvironment, c.Environment.Kind,
		EnvConfigServerHome, c.Home,
	)
}
