// Package config manges all the ConfigServer configuration activities
package config

import (
	"fmt"
	"go.uber.org/zap"
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

	var envHome = os.Getenv(EnvConfigServerHome)
	var envCfg = os.Getenv(EnvConfigServerCfg)
	var envRepos = os.Getenv(EnvRepositoriesHome)

	root := configurationRoot
	if len(root) == 0 {
		root = envCfg
		if len(root) == 0 {
			root = envHome
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

	if len(root) == 0 {
		root = DefaultHome
	}
	homeFolder, err := os.Stat(envHome)
	if err != nil {
		return nil, fmt.Errorf("cannot stat '%s': %w", root, err)
	}
	if !homeFolder.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", envHome)
	}

	if len(envRepos) == 0 {
		envRepos = path.Join(root, "repositories")
	}
	repoFolder, err := os.Stat(envRepos)
	if err != nil {
		return nil, fmt.Errorf("cannot stat '%s': %s", root, err.Error())
	}
	if !repoFolder.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", envRepos)
	}

	zap.L().Sugar().Infof("%s='%s'", EnvConfigServerHome, envHome)
	zap.L().Sugar().Infof("%s='%s'", EnvConfigServerCfg, envCfg)
	zap.L().Sugar().Infof("%s='%s'", EnvRepositoriesHome, envRepos)

	v := viper.New()
	v.AddConfigPath(root)
	v.SetConfigName("configserver")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading configserver.yaml from '%s' %w", root, err)
	}

	conf := &Config{
		ListenOn:                     ":8090",
		CacheEvictorIntervalSeconds:  10,
		CacheStorageSeconds:          30,
		Home:                         envHome,
		RepositoriesCheckoutLocation: envRepos,
		LoadedFrom:                   path.Join(root, "configserver.yaml"),
	}
	err = v.Unmarshal(&conf)
	if err != nil {
		return nil, fmt.Errorf("error reading configserver.yaml from '%s': %w", root, err)
	}
	return conf, nil
}

// Config stores all the supported configuration options for a ConfigServer Instance
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
