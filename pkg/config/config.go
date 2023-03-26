package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const ENV_CONFIGSERVER_HOME = "CONFIGSERVER_HOME"
const DEFAULT_HOME string = "/var/run/configserver"

// Reads the configuration from the location pointed by the provided configurationRoot parameter.
// If this parameter is omited, tries to locate the configuration using the $CONFIGSERVER_HOME env variable.
// If the variable is not defined, uses /var/run/configserver as a default.
func Read(configurationRoot string) (*Config, error) {
	root := configurationRoot
	if len(root) == 0 {
		root = os.Getenv(ENV_CONFIGSERVER_HOME)
		if len(root) == 0 {
			root = DEFAULT_HOME
		}
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

	conf := &Config{}
	v.Unmarshal(&conf)
	conf.Source = path.Join(root, "configserver.yaml")
	return conf, nil
}

type Config struct {
	Key          string
	Source       string
	Repositories Repositories
}
