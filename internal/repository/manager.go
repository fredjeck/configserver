package repository

import (
	"fmt"
	"github.com/fredjeck/configserver/internal/config"
	"log/slog"
)

// Manager is one-stop shop for managing multiple repositories configured via yaml files
type Manager struct {
	Configuration *config.Repositories   // git configuration
	Repositories  map[string]*Repository // list of configured repository
	Heartbeat     chan UpdateEvent       // uplink channel used by beholders to communicate
}

// NewManager creates a new repository manager by parsing the provided target repository configuration location
func NewManager(configuration *config.Repositories) (*Manager, error) {

	hb := make(chan UpdateEvent)

	repos := make(map[string]*Repository)
	for _, repo := range configuration.Configuration {

		if len(repo.Branch) == 0 {
			repo.Branch = "main"
		}
		repos[repo.Name] = &Repository{
			Configuration: repo,
			Beholder:      NewBeholder(configuration.CheckoutLocation, repo, hb),
			Statistics:    &Statistics{},
		}
	}

	return &Manager{configuration, repos, hb}, nil
}

// Start will generate a repository beholder for each found configuration and will attempt to create a local copy
func (mgr *Manager) Start() {
	go mgr.listen()
	for name, repo := range mgr.Repositories {
		slog.Info("starting beholder", "name", name)
		repo.Beholder.Watch()
	}
}

// Get scans the target repository for the file pointed by the provided path
func (mgr *Manager) Get(repository string, path string) ([]byte, error) {
	r, ok := mgr.Repositories[repository]
	if !ok {
		return nil, fmt.Errorf("'%s' is not a valid repository", repository)
	}

	if !r.Beholder.Active {
		return nil, fmt.Errorf("'%s' cannot be checked out due to %w", repository, r.Statistics.LastError)
	}

	contents, err := r.Beholder.File(path)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// IsClientAllowed checks if the provided client identifier is configured for repository access
func (mgr *Manager) IsClientAllowed(repository string, client string) bool {
	if repo, ok := mgr.Repositories[repository]; ok {
		if repo.IsClientAllowed(client) {
			return true
		}
	}
	return false
}

// listen reads the heartbeat channel for beholder events
func (mgr *Manager) listen() {
	for event := range mgr.Heartbeat {
		mgr.Repositories[event.RepositoryName].Statistics.LastError = event.LastError
		mgr.Repositories[event.RepositoryName].Statistics.NextUpdate = event.NextUpdate
		mgr.Repositories[event.RepositoryName].Statistics.LastUpdate = event.LastUpdate
	}
}
