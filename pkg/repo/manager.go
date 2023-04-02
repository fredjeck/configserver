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

func NewRepositoryManger(conf config.Config, logger zap.Logger) *RepositoryManager {
	return &RepositoryManager{configuration: conf, repositoriesRoot: path.Join(conf.Home, REPOSITORIES), logger: logger}
}

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
