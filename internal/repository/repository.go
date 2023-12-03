package repository

import (
	"fmt"
	"github.com/fredjeck/configserver/internal/config"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"
)

// Configuration represents a repository configuration as retrieved from yaml configuration files
type Configuration struct {
	Name                   string `yaml:"name"`
	Url                    string `yaml:"url"`
	Branch                 string `yaml:"branch"`
	RefreshIntervalSeconds int    `yaml:"refreshIntervalSeconds"`
	CheckoutLocation       string `yaml:"checkoutLocation"`
	Token                  string `yaml:"token"`
}

type Statistics struct {
	Hitcount   int64     `json:"hitCount"`
	LastUpdate time.Time `json:"lastUpdate"`
	NextUpdate time.Time `json:"nextUpdate"`
	LastError  error     `json:"lastError"`
}

type UpdateEvent struct {
	RepositoryName string
	LastUpdate     time.Time
	NextUpdate     time.Time
	LastError      error
	Active         bool
}

type Repository struct {
	Configuration *Configuration
	Statistics    *Statistics
	Beholder      *Beholder
	Active        bool
}

// Manager is one-stop shop for managing all the configured repositories
type Manager struct {
	*config.GitConfiguration
	Repositories map[string]*Repository
	Heartbeat    chan UpdateEvent
}

func NewManager(configuration *config.GitConfiguration) (*Manager, error) {

	hb := make(chan UpdateEvent)
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

		repo := &Configuration{}
		err = yaml.Unmarshal([]byte(data), &repo)
		if err != nil {
			return nil, fmt.Errorf("'%s' repository configuration cannot be loaded : %w", repoConfigPath, err)
		}
		if len(repo.Branch) == 0 {
			repo.Branch = "main"
		}
		repos[repo.Name] = &Repository{
			Configuration: repo,
			Beholder:      NewBeholder(repo, hb),
			Statistics:    &Statistics{},
			Active:        true,
		}
		slog.Info("repository configuration loaded from", "location", repoConfigPath)
	}

	return &Manager{configuration, repos, hb}, nil
}

func (mgr Manager) Start() {
	for name, repo := range mgr.Repositories {
		slog.Info("Starting beholder", "name", name)
		repo.Beholder.Watch()
	}
}

func (mgr *Manager) listen() {
	for event := range mgr.Heartbeat {

		if event.LastError != nil {
			slog.Error("error", event.LastError)
		}

		mgr.Repositories[event.RepositoryName].Statistics.LastError = event.LastError
		mgr.Repositories[event.RepositoryName].Statistics.NextUpdate = event.NextUpdate
		mgr.Repositories[event.RepositoryName].Statistics.LastUpdate = event.LastUpdate
		mgr.Repositories[event.RepositoryName].Active = event.Active
	}
}
