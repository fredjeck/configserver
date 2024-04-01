// Package configuration groups all configuration related utilities
package configuration

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// EnvConfigServerEnvironment holds the name of the environment variable storing the current active environment type (dev, int, prod ...)
const EnvConfigServerEnvironment string = "CONFIGSERVER_ENV"

// EnvConfigServerHome holds the path to the directory containing the application's config
const EnvConfigServerHome string = "CONFIGSERVER_HOME"

// Configuration groups all the supported configuration options
type Configuration struct {
	Source string
	// Environment variables for ease of use
	*Environment
	// Server related settings
	*Server
	// Repositories configuration
	*Repositories
}

// Environment gathers all the environment variable used by ConfigServer
type Environment struct {
	Kind string // environment kind i,e dev, int, prod
	Home string // path to the directory containing configserver configuration files
}

// Server groups all the configserver related settings
type Server struct {
	PassPhrase             string `yaml:"passPhrase"`             // key used for secret encryption
	ListenOn               string `yaml:"listenOn"`               // address and port on which the server will listen for incoming requests
	SecretExpiryDays       int    `yaml:"secretExpiryDays"`       // number of days after which a secret is considered as expired
	ValidateSecretLifeSpan bool   `yaml:"validateSecretLifespan"` // if true, an expired secret will be considered invalid
}

// Repositories materializes the GIT repositories configuration
type Repositories struct {
	CheckoutLocation string        `yaml:"checkoutLocation"` // Folder to which the repositories are stored
	Configuration    []*Repository `yaml:"configuration"`    // Collection of git repositories configuration
}

// Repository is a single GIT repository configuration
type Repository struct {
	Name                   string   `yaml:"name"`
	URL                    string   `yaml:"url"`
	Branch                 string   `yaml:"branch"`
	RefreshIntervalSeconds int      `yaml:"refreshIntervalSeconds"`
	CheckoutLocation       string   `yaml:"checkoutLocation"`
	Token                  string   `yaml:"token"`
	Clients                []string `yaml:"clients"`
}

// DefaultConfiguration for when its needed
var DefaultConfiguration = &Configuration{
	Environment: &Environment{
		Kind: "production",
		Home: "/var/run/configserver",
	},
	Server: &Server{
		PassPhrase:             "This is a default passphrase and should be changed",
		ListenOn:               ":4200",
		SecretExpiryDays:       365,
		ValidateSecretLifeSpan: false,
	},
	Repositories: &Repositories{
		CheckoutLocation: "",
	},
}

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

	home := os.Getenv(EnvConfigServerHome)
	if len(home) == 0 {
		home = "/var/run/configserver"
	}

	configPath := path
	if _, err := os.Stat(path); err != nil {
		configPath = filepath.Join(home, "configserver.yml")
		if _, err := os.Stat(configPath); err != nil {
			return nil, fmt.Errorf("'%s' configserver configuration cannot be found or is not accessible: %w", configPath, err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("'%s' configserver configuration cannot be loaded : %w", configPath, err)
	}

	config := DefaultConfiguration
	config.Source = configPath

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("'%s' cannot unmarshal yaml file : %w", path, err)
	}

	if config.Repositories.CheckoutLocation == "" {
		config.Repositories.CheckoutLocation = filepath.Join(home, "repositories")
		slog.Info(fmt.Sprintf("Repositories checkout location defaulted to '%s'", config.Repositories.CheckoutLocation))
	}

	return config, nil
}

// LogEnvironment logs the current environment configuration
func (c *Configuration) LogEnvironment() {
	slog.Info("Configserver Runtime Environment",
		EnvConfigServerEnvironment, c.Environment.Kind,
		EnvConfigServerHome, c.Home,
		"configuration.source", c.Source,
	)
}
