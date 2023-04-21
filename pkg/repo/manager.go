package repo

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fredjeck/configserver/pkg/config"
	"go.uber.org/zap"
)

type RepositoryManager struct {
	repositoriesRoot string
	logger           zap.Logger
	Repositories     map[string]*RepositoryEntry
	heartbeat        chan RepositoryUpdateEvent
}

type RepositoryEntry struct {
	Hitcount      int64
	LastUpdate    time.Time
	NextUpdate    time.Time
	Configuration config.Repository
}

type RepositoryUpdateEvent struct {
	name       string
	lastUpdate time.Time
	nextUpdate time.Time
}

const RepositoriesFolder string = "repositories"

var (
	ErrRepositoryNotFound = errors.New("repository does not exist")
	ErrFileNotFound       = errors.New("file not found in repository")
	ErrInvalidPath        = errors.New("unsupported repository path")
)

// NewManager creates a new repository manager
func NewManager(conf config.Config, logger zap.Logger) *RepositoryManager {
	r := make(map[string]*RepositoryEntry)
	for _, v := range conf.Repositories {
		r[v.Name] = &RepositoryEntry{
			Configuration: v,
			Hitcount:      0,
			LastUpdate:    time.Time{},
			NextUpdate:    time.Time{},
		}
	}

	return &RepositoryManager{
		Repositories:     r,
		repositoriesRoot: path.Join(conf.Home, RepositoriesFolder),
		logger:           logger,
		heartbeat:        make(chan RepositoryUpdateEvent),
	}
}

// Checkout starts asynchronously checking out the configured repositories
// and maintains their local copy up to date based on the repository configuration
func (mgr *RepositoryManager) Checkout() error {
	if _, err := os.Stat(mgr.repositoriesRoot); os.IsNotExist(err) {
		err = os.MkdirAll(mgr.repositoriesRoot, 0700)
		if err != nil {
			return err
		}
	}

	for _, repository := range mgr.Repositories {
		repositoryPath := path.Join(mgr.repositoriesRoot, repository.Configuration.Name)
		go NewWatcher(repository.Configuration, repositoryPath, mgr.logger, mgr.heartbeat)
	}
	go mgr.ReadUpdates()
	return nil
}

func (mgr *RepositoryManager) ReadUpdates() {
	for event := range mgr.heartbeat {
		mgr.Repositories[event.name].NextUpdate = event.nextUpdate
		mgr.Repositories[event.name].LastUpdate = event.lastUpdate
	}
}

// Get retrieves the file from repository at the specified path
func (mgr *RepositoryManager) Get(repository string, target string) ([]byte, error) {

	r, found := mgr.Repositories[strings.ToLower(repository)]
	if !found {
		return nil, ErrRepositoryNotFound
	}

	r.Hitcount++

	if strings.Contains(target, "."+string(os.PathSeparator)) || strings.Contains(target, ".."+string(os.PathSeparator)) {
		return nil, ErrInvalidPath
	}

	// TODO sanitize to avoid browsing filesystem
	repositoryPath := path.Join(mgr.repositoriesRoot, repository, target)

	content, err := os.ReadFile(repositoryPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrFileNotFound
	}

	return content, nil
}
