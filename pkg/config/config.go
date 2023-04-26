package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	// EnvConfigserverHome defines the name of the environment variable pointing to the configuration file
	EnvConfigServerHome = "CONFIGSERVER_HOME"
	// EnvConfigServerCfg defines the main configuration file
	EnvConfigServerCfg = "CONFIGSERVER_CFG"
	// EnvRepositoriesHome defines the name of the environment variable pointing to the location where repositories will be checked out
	EnvRepositoriesHome = "CONFIGSERVER_REPOSITORIES"
	// DefaultHome defines the default home directory used when the CONFIGSERVER_HOME environment variable is not defined
	DefaultHome string = "/var/run/configserver"
)

// ReadFromPath loads the configuration from the location pointed by the provided configurationRoot parameter.
// If this parameter is omitted, tries to locate the configuration using the $CONFIGSERVER_HOME env variable.
// If the variable is not defined, uses /var/run/configserver as a default.
func ReadFromPath(configurationRoot string) (*Config, error) {
	root := configurationRoot
	if len(root) == 0 {
		root = os.Getenv(EnvConfigServerCfg)
		if len(root) == 0 {
			root = os.Getenv(EnvConfigServerHome)
			if len(root) == 0 {
				root = DefaultHome
			}
		}
	}
	cfgRoot, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot stat '%s': %s", root, err.Error())
	}
	if !cfgRoot.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", root)
	}

	home := os.Getenv(EnvConfigServerHome)
	if len(root) == 0 {
		root = DefaultHome
	}
	homeFolder, err := os.Stat(home)
	if err != nil {
		return nil, fmt.Errorf("cannot stat '%s': %s", root, err.Error())
	}
	if !homeFolder.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", home)
	}

	repositoriesLocation := os.Getenv(EnvRepositoriesHome)
	if len(repositoriesLocation) == 0 {
		repositoriesLocation = path.Join(root, "repositories")
	}
	repoFolder, err := os.Stat(repositoriesLocation)
	if err != nil {
		return nil, fmt.Errorf("cannot stat '%s': %s", root, err.Error())
	}
	if !repoFolder.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", repositoriesLocation)
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
		Home:                         home,
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
