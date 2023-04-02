package repo

import (
	"os"
	"path"

	"github.com/fredjeck/configserver/pkg/config"
	"go.uber.org/zap"
)

type RepositoryManager struct {
	configuration    config.Config
	repositoriesRoot string
	logger           zap.Logger
}

const REPOSITORIES string = "repositories"

// Creates a new repository manager
func NewManager(conf config.Config, logger zap.Logger) *RepositoryManager {
	return &RepositoryManager{configuration: conf, repositoriesRoot: path.Join(conf.Home, REPOSITORIES), logger: logger}
}

// Starts asynchroneously checking out the configured repositories
// and maintains their local copy up to date based on the repository configuration
func (mgr *RepositoryManager) Checkout() error {
	if _, err := os.Stat(mgr.repositoriesRoot); os.IsNotExist(err) {
		err = os.MkdirAll(mgr.repositoriesRoot, 0700)
		if err != nil {
			return err
		}
	}

	for _, repository := range mgr.configuration.Repositories {
		repositoryPath := path.Join(mgr.repositoriesRoot, repository.Name)
		go Watch(repository, repositoryPath, mgr.logger)
	}
	return nil
}

// Gets the file from repository at the specified path
func (mgr RepositoryManager) Get(repository string, path string) ([]byte, error) {
	return nil, nil
}
