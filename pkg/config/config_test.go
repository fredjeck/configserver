package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assetsPath() string {
	wd, _ := os.Getwd()
	source := path.Join(wd, "..", "..", "test", "config")
	return source
}

func TestReadFromPath(t *testing.T) {
	source := path.Join(assetsPath(), "valid")
	config, err := Read(source)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, path.Join(source, "configserver.yaml"), config.Source)
}

func TestReadFromEnvironment(t *testing.T) {
	assets := path.Join(assetsPath(), "valid")
	os.Setenv(ENV_CONFIGSERVER_HOME, assets)
	var source string
	config, err := Read(source)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, path.Join(assets, "configserver.yaml"), config.Source)
}
