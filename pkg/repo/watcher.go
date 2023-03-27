package repo

import (
	"os"
	"time"

	"github.com/fredjeck/configserver/pkg/config"
	"github.com/go-git/go-git/v5"
)

func Watch(repository config.Repository, localPath string) {

	for {
		_, err := git.PlainClone(localPath, false, &git.CloneOptions{
			URL:      repository.Url,
			Progress: os.Stdout,
		})

		if err != nil {
			return
		}

		time.Sleep(time.Duration(repository.RefreshInterval) * time.Second)
	}
}
