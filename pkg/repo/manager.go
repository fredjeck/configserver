package repo

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/fredjeck/configserver/pkg/config"
	"go.uber.org/zap"
)

type RepositoryManager struct {
	configuration    config.Config
	repositoriesRoot string
	logger           zap.Logger
}

const REPOSITORIES string = "repositories"

var (
	ErrRepositoryNotFound = errors.New("repository does not exist")
	ErrFileNotFound       = errors.New("file not found in repository")
	ErrInvalidPath        = errors.New("unsupported repository path")
)

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
func (mgr RepositoryManager) Get(repository string, target string) ([]byte, error) {
	found := false
	for r := range mgr.configuration.Repositories {
		if strings.EqualFold(mgr.configuration.Repositories[r].Name, repository) {
			found = true
		}
	}
	if !found {
		return nil, ErrRepositoryNotFound
	}

	if strings.Contains(target, "."+string(os.PathSeparator)) || strings.Contains(target, ".."+string(os.PathSeparator)) {
		return nil, ErrInvalidPath
	}

	// TODO sanitize to avoid browsing filesystem
	repositoryPath := path.Join(mgr.repositoriesRoot, repository, target)

	content, error := os.ReadFile(repositoryPath)
	if errors.Is(error, os.ErrNotExist) {
		return nil, ErrFileNotFound
	}

	return content, nil

}
