package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	// EnvConfigserverHome defines the name of the environment variable pointing to the configuration file
	EnvConfigserverHome = "CONFIGSERVER_HOME"
	// EnvRepositoriesHome defines the name of the environment variable pointing to the location where repositories will be checked out
	EnvRepositoriesHome = "CONFIGSERVER_REPOSITORIES"
	// DefaultHome defines the default home directory used when the CONFIGSERVER_HOME environment variable is not defined
	DefaultHome string = "/var/run/configserver"
)

// Read loads the configuration from the location pointed by the $CONFIGSERVER_HOME env variable.
// If the variable is not defined, uses /var/run/configserver as a default.
func Read() (*Config, error) {
	root := os.Getenv(EnvConfigserverHome)
	if len(root) == 0 {
		root = DefaultHome
	}
	return ReadFromPath(root)
}

// ReadFromPath loads the configuration from the location pointed by the provided configurationRoot parameter.
// If this parameter is omitted, tries to locate the configuration using the $CONFIGSERVER_HOME env variable.
// If the variable is not defined, uses /var/run/configserver as a default.
func ReadFromPath(configurationRoot string) (*Config, error) {
	root := configurationRoot
	if len(root) == 0 {
		root = os.Getenv(EnvConfigserverHome)
		if len(root) == 0 {
			root = DefaultHome
		}
	}

	repositoriesLocation := os.Getenv(EnvRepositoriesHome)
	if len(repositoriesLocation) == 0 {
		repositoriesLocation = path.Join(root, "repositories")
	}

	home, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot stat '%s': %s", root, err.Error())
	}

	if !home.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", root)
	}

	v := viper.New()
	v.AddConfigPath(root)
	v.SetConfigName("configserver")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading configserver.yaml from '%s': %s", root, err.Error())
	}

	conf := &Config{
		ListenOn:                     ":8090",
		CacheEvictorIntervalSeconds:  10,
		CacheStorageSeconds:          30,
		Home:                         root,
		RepositoriesCheckoutLocation: repositoriesLocation,
		LoadedFrom:                   path.Join(root, "configserver.yaml"),
	}
	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("error reading configserver.yaml from '%s': %s", root, err.Error())
	}
	return conf, nil
}

type Config struct {
	ListenOn                     string
	CacheEvictorIntervalSeconds  int
	CacheStorageSeconds          int
	LoadedFrom                   string
	Home                         string
	RepositoriesCheckoutLocation string
	Repositories                 Repositories
}

// EncryptionKeyPath returns the path to the location of the encryption.key file
func (config Config) EncryptionKeyPath() string {
	return path.Join(config.Home, "encryption.key")
}
