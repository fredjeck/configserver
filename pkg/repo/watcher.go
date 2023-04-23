package repo

import (
	"errors"
	"fmt"
	"github.com/fredjeck/configserver/pkg/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
	"os"
	"time"
)

type RepositoryWatcher struct {
	repository config.Repository
	localPath  string
	logger     zap.Logger
	uplink     chan RepositoryUpdateEvent
}

// NewWatcher uses the provided repository configuration to periodically pulls the repository
// to the path pointed by localPath.
// If and error occurs, provides a detailed output in the logs
func NewWatcher(repo config.Repository, localPath string, logger zap.Logger, uplink chan RepositoryUpdateEvent) *RepositoryWatcher {
	return &RepositoryWatcher{
		repository: repo,
		localPath:  localPath,
		logger:     logger,
		uplink:     uplink,
	}
}

func (w *RepositoryWatcher) Watch() {
	go w.watchInternal()
}

func (w *RepositoryWatcher) broadcast(last time.Time, next time.Time, error string) {
	w.uplink <- *&RepositoryUpdateEvent{
		lastUpdate:     last,
		nextUpdate:     next,
		lastError:      error,
		repositoryName: w.repository.Name,
	}
}

func (w *RepositoryWatcher) watchInternal() {
	log := w.logger.With(zap.String("repository.name", w.repository.Name)).With(zap.String("repository.url", w.repository.Url)).With(zap.String("repository.localPath", w.localPath))
	for {
		last := time.Now()
		next := time.Now()
		log.Sugar().Infof("pulling repository '%s'", w.repository.Name)

		workspace, err := git.PlainOpen(w.localPath)
		if err != nil {
			log.Sugar().Infof("no local copy of '%s' was found... cloning", w.repository.Name)
			err := os.MkdirAll(w.localPath, 0700)
			if err != nil {
				msg := fmt.Sprintf("an error occured while preparing %s for checkout: %s", w.localPath, err.Error())
				w.broadcast(last, next, msg)
				log.Sugar().Error(msg)
				return
			}
			_, err = git.PlainClone(w.localPath, false, &git.CloneOptions{
				URL:      w.repository.Url,
				Progress: os.Stdout,
				Auth:     &http.BasicAuth{Username: "user", Password: w.repository.Token},
			})
			if err != nil {
				msg := fmt.Sprintf("could not clone '%s' to '%s' : %s", w.repository.Url, w.localPath, err.Error())
				w.broadcast(last, next, msg)
				log.Sugar().Error(msg)
				return
			}
		} else {
			var tree *git.Worktree
			tree, err = workspace.Worktree()
			if err != nil {
				msg := fmt.Sprintf("'%s' : unable to open local copy : %s", w.localPath, err.Error())
				w.broadcast(last, next, msg)
				log.Sugar().Error(msg)
				return
			}

			err = tree.Pull(&git.PullOptions{
				RemoteName: "origin",
				Auth:       &http.BasicAuth{Username: "user", Password: w.repository.Token},
				Force:      true,
			})
			if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
				msg := fmt.Sprintf("'%s' : unable to pull latest changes : %s", w.repository.Name, err.Error())
				w.broadcast(last, next, msg)
				log.Sugar().Error(msg)
				return
			}
		}

		nextRefresh := time.Duration(w.repository.RefreshInterval) * time.Second
		log.Sugar().Infof("'%s' : next pull will occur @ %s", w.repository.Name, time.Now().Add(nextRefresh))
		w.broadcast(last, time.Now().Add(nextRefresh), "")
		time.Sleep(nextRefresh)
	}
}
