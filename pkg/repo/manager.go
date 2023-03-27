package repo

import (
	"os"
	"path"

	"github.com/fredjeck/configserver/pkg/config"
)

type RepositoryManager struct {
	configuration    config.Config
	repositoriesRoot string
}

const REPOSITORIES string = "repositories"

func NewRepositoryManger(conf config.Config) *RepositoryManager {
	return &RepositoryManager{configuration: conf, repositoriesRoot: path.Join(conf.Home, REPOSITORIES)}
}

func (mgr *RepositoryManager) Checkout() error {
	if _, err := os.Stat(mgr.repositoriesRoot); os.IsNotExist(err) {
		err = os.MkdirAll(mgr.repositoriesRoot, 0644)
		if err != nil {
			return err
		}
	}

	for _, repository := range mgr.configuration.Repositories {
		go Watch(repository, mgr.repositoriesRoot)
	}
	return nil
}
