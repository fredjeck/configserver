package repository

import (
	"fmt"
	"github.com/fredjeck/configserver/internal/config"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// dummy code

func dummy() {
	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git",
		Progress: os.Stdout,
	})

	CheckIfError(err)
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Repository represents a repository configuration as retrieved from yaml configuration files
type Repository struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

// Manager is one stop shop for managing all the configured repositories
type Manager struct {
	*config.GitConfiguration
	Repositories map[string]*Repository
}

func NewManager(configuration *config.GitConfiguration) (*Manager, error) {

	entries, err := os.ReadDir(configuration.RepositoriesConfigurationLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to parse '%s' for repository configurations: %w", configuration.RepositoriesConfigurationLocation, err)
	}

	repos := make(map[string]*Repository)
	for _, element := range entries {
		if element.IsDir() || (".yaml" != filepath.Ext(element.Name()) && ".yml" != filepath.Ext(element.Name())) {
			continue
		}

		repoConfigPath := path.Join(configuration.RepositoriesConfigurationLocation, element.Name())
		data, err := os.ReadFile(repoConfigPath)
		if err != nil {
			return nil, fmt.Errorf("'%s' repository configuration cannot be loaded : %w", repoConfigPath, err)
		}

		repo := &Repository{}
		err = yaml.Unmarshal([]byte(data), &repo)
		if err != nil {
			return nil, fmt.Errorf("'%s' repository configuration cannot be loaded : %w", repoConfigPath, err)
		}
		repos[repo.Name] = repo
		slog.Info("Repository configuration loaded from", "location", repoConfigPath)
	}

	return &Manager{configuration, repos}, nil
}
