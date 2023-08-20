package config

import (
	"go.uber.org/zap"
	"os"
	"strings"
)

func InitLogging() {
	env := strings.ToLower(os.Getenv(ConfigServerEnvironment))
	if strings.Contains(env, "dev") {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	}
	zap.L().Sugar().Infof("%s = '%s'", ConfigServerEnvironment, strings.ToUpper(env))
}

func LoadFrom(path string) (*Configuration, error) {
	return &Configuration{}, nil
}

type Configuration struct {
}
