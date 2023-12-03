package repository

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log/slog"
	"os"
	"time"
)

type Beholder struct {
	*Configuration
	heartbeat chan UpdateEvent
}

func NewBeholder(configuration *Configuration, heartbeat chan UpdateEvent) *Beholder {
	return &Beholder{configuration, heartbeat}
}

func (w *Beholder) Watch() {
	go w.watchInternal()
}

func (w *Beholder) watchInternal() {
	for {
		last := time.Now()
		next := time.Now()
		slog.Info("pulling repository", "repository_name", w.Name)

		workspace, err := git.PlainOpen(w.CheckoutLocation)
		if err != nil {
			slog.Info("no local copy found, creating a fresh clone", "repository_name", w.Name)
			err := os.MkdirAll(w.CheckoutLocation, 0700)
			if err != nil {
				e := fmt.Errorf("an error occured while preparing %s for checkout: %w", w.CheckoutLocation, err)
				w.broadcast(last, next, e, false)
				return
			}
			_, err = git.PlainClone(w.CheckoutLocation, false, &git.CloneOptions{
				URL:        w.Url,
				Progress:   os.Stdout,
				Auth:       &http.BasicAuth{Username: "user", Password: w.Token},
				RemoteName: w.Branch,
			})
			if err != nil {
				e := fmt.Errorf("could not clone '%s' to '%s' : %w", w.Url, w.CheckoutLocation, err)
				w.broadcast(last, next, e, false)
				return
			}
		} else {
			var tree *git.Worktree
			tree, err = workspace.Worktree()
			if err != nil {
				e := fmt.Errorf("'%s' : unable to open local copy : %s", w.CheckoutLocation, err)
				w.broadcast(last, next, e, false)
				return
			}

			err = tree.Pull(&git.PullOptions{
				RemoteName: "origin",
				Auth:       &http.BasicAuth{Username: "user", Password: w.Token},
				Force:      true,
			})

			if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
				e := fmt.Errorf("'%s' : unable to pull latest changes : %s", w.Name, err)
				w.broadcast(last, next, e, false)
				return
			}
		}

		nextRefresh := time.Duration(w.RefreshIntervalSeconds) * time.Second
		slog.Info(fmt.Sprintf("'%s' next pull will occur @ %s", w.Name, time.Now().Add(nextRefresh)), "repository_name", w.Name)
		w.broadcast(last, time.Now().Add(nextRefresh), nil, true)
		time.Sleep(nextRefresh)
	}
}

func (w Beholder) broadcast(last time.Time, next time.Time, error error, active bool) {
	w.heartbeat <- *&UpdateEvent{
		LastUpdate:     last,
		NextUpdate:     next,
		LastError:      error,
		RepositoryName: w.Name,
		Active:         active,
	}
}
