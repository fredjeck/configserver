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
	repositories     config.Repositories
	repositoriesRoot string
	logger           zap.Logger
}

const RepositoriesFolder string = "repositories"

var (
	ErrRepositoryNotFound = errors.New("repository does not exist")
	ErrFileNotFound       = errors.New("file not found in repository")
	ErrInvalidPath        = errors.New("unsupported repository path")
)

// NewManager creates a new repository manager
func NewManager(conf config.Config, logger zap.Logger) *RepositoryManager {
	return &RepositoryManager{repositories: conf.Repositories, repositoriesRoot: path.Join(conf.Home, RepositoriesFolder), logger: logger}
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

	for _, repository := range mgr.repositories {
		repositoryPath := path.Join(mgr.repositoriesRoot, repository.Name)
		go NewWatcher(repository, repositoryPath, mgr.logger)
	}
	return nil
}

// Get retrieves the file from repository at the specified path
func (mgr *RepositoryManager) Get(repository string, target string) ([]byte, error) {
	found := false
	for r := range mgr.repositories {
		if strings.EqualFold(mgr.repositories[r].Name, repository) {
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

	content, err := os.ReadFile(repositoryPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrFileNotFound
	}

	return content, nil
}
