package repo

import (
	"errors"
	"os"
	"time"

	"github.com/fredjeck/configserver/pkg/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
)

// Using the provided repository configuration periodically pulls the repository
// to the path pointed by localPath.
// If and error occurs, provides a detailed output in the logs
func Watch(repository config.Repository, localPath string, logger zap.Logger) {

	log := logger.With(zap.String("repository.name", repository.Name)).With(zap.String("repository.url", repository.Url)).With(zap.String("repository.localPath", localPath))
	for {
		log.Sugar().Infof("pulling repository '%s'", repository.Name)

		workspace, err := git.PlainOpen(localPath)
		if err != nil {
			log.Sugar().Infof("no local copy of '%s' was found... cloning", repository.Name)
			os.MkdirAll(localPath, 0700)
			_, err = git.PlainClone(localPath, false, &git.CloneOptions{
				URL:      repository.Url,
				Progress: os.Stdout,
				Auth:     &http.BasicAuth{Username: "user", Password: repository.Token},
			})
			if err != nil {
				log.Sugar().Errorf("could not clone '%s' to '%s' : %s", repository.Url, localPath, err.Error())
				return
			}
		} else {
			var tree *git.Worktree
			tree, err = workspace.Worktree()
			if err != nil {
				log.Sugar().Errorf("'%s' : unable to open local copy : %s", localPath, err.Error())
				return
			}

			err = tree.Pull(&git.PullOptions{
				RemoteName: "origin",
				Auth:       &http.BasicAuth{Username: "user", Password: repository.Token},
				Force:      true,
			})
			if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
				log.Sugar().Errorf("'%s' : unable to pull latest changes : %s", repository.Name, err.Error())
				return
			}
		}

		next := time.Duration(repository.RefreshInterval) * time.Second
		log.Sugar().Infof("'%s' : next pull will occur @ %s", repository.Name, time.Now().Add(next))
		time.Sleep(next)
	}
}
