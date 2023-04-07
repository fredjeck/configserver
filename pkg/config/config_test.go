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
	config, err := ReadFromPath(source)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, path.Join(source, "configserver.yaml"), config.LoadedFrom)
}

func TestReadFromEnvironment(t *testing.T) {
	assets := path.Join(assetsPath(), "valid")
	_ = os.Setenv(EnvConfigserverHome, assets)
	var source string
	config, err := ReadFromPath(source)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, path.Join(assets, "configserver.yaml"), config.LoadedFrom)
}

func TestReadFromFile(t *testing.T) {
	assets := path.Join(assetsPath(), "valid", "configserver.yaml")
	_, err := ReadFromPath(assets)
	if err == nil {
		t.Error(err)
	}
}

func TestReadFromUnexistingLocation(t *testing.T) {
	assets := path.Join("invalidlocation")
	_, err := ReadFromPath(assets)
	if err == nil {
		t.Error(err)
	}
}
