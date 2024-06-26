package repository

import "C"
import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/fredjeck/configserver/internal/configuration"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/uuid"
)

// Beholder are responsible for maintaining local copies of git repositories up to date based on the provided configuration
// As they are running in background they make use of a heartbeat channel towards their initiator to communicate about
// repositories update event
type Beholder struct {
	configuration    *configuration.Repository // The current configured repository
	checkoutLocation string                    // Place where the repositories are checked out
	Active           bool                      // true if the repository is actively monitored and can be used
	heartbeat        chan UpdateEvent          // Uplink to the beholder's initiator
	mutex            *sync.RWMutex             // Used to ensure no read operation is allowed while the repository is being updated
}

// NewBeholder initiates a new beholder for the provided configuration
// a call to Watch() is mandatory to start the beholder process
func NewBeholder(checkoutLocation string, configuration *configuration.Repository, heartbeat chan UpdateEvent) *Beholder {
	uid := uuid.New()
	return &Beholder{configuration, filepath.Join(checkoutLocation, uid.String()), true, heartbeat, &sync.RWMutex{}}
}

// Watch initiates the creation of a local copy of the configured repository and will periodically update the repository
// to its latest state
func (w *Beholder) Watch() {
	go w.watchInternal()
}

// see Watch
func (w *Beholder) watchInternal() {
	var lastError error

	for {
		w.mutex.Lock()
		last := time.Now()

		slog.Info("cloning repository", logKeyRepositoryName, w.configuration.Name, logKeyCheckoutLocation, w.checkoutLocation, logKeyRepositoryURL, w.configuration.URL)

		if err := os.MkdirAll(w.checkoutLocation, os.ModePerm); err != nil {
			lastError = fmt.Errorf("cannot create path '%s' to checkout '%s': %w", w.checkoutLocation, w.configuration.Name, err)
			break
		}

		workspace, err := git.PlainOpen(w.checkoutLocation)
		if err != nil {
			slog.Info("no local copy found, creating a fresh clone", logKeyRepositoryName, w.configuration.Name)
			workspace, err = git.PlainClone(w.checkoutLocation, false, &git.CloneOptions{
				URL:      w.configuration.URL,
				Progress: os.Stdout,
			})

			if err != nil {
				lastError = fmt.Errorf("could not clone '%s' to '%s' : %w", w.configuration.URL, w.checkoutLocation, err)
				break
			}
		}

		tree, err := workspace.Worktree()
		if err != nil {
			lastError = fmt.Errorf("'%s' : unable to open local copy : %w", w.checkoutLocation, err)
			break
		}

		if len(w.configuration.Branch) > 0 {
			// Fetch remote branches
			err = workspace.Fetch(&git.FetchOptions{
				RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
			})
			if err != nil {
				lastError = fmt.Errorf("'%s' : unable to fetch repository : %w", w.configuration.URL, err)
				break
			}
			err = tree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(w.configuration.Branch),
				Force:  true,
			})
			if err != nil {
				lastError = fmt.Errorf("'%s' : unable to checkout branch '%s': %w", w.configuration.URL, w.configuration.Branch, err)
				break
			}
		}

		err = tree.Pull(&git.PullOptions{
			Force: true,
		})

		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			lastError = fmt.Errorf("'%s' : unable to pull latest changes : %s", w.configuration.Name, err)
			break
		}

		nextRefresh := time.Duration(w.configuration.RefreshIntervalSeconds) * time.Second
		slog.Info(fmt.Sprintf("'%s' next pull will occur @ %s", w.configuration.Name, time.Now().Add(nextRefresh)), logKeyRepositoryName, w.configuration.Name, logKeyCheckoutLocation, w.checkoutLocation, logKeyRepositoryURL, w.configuration.URL)
		w.broadcast(last, time.Now().Add(nextRefresh), nil)
		w.mutex.Unlock()
		time.Sleep(nextRefresh)
	}

	// If we are here something bad happend
	slog.Error("Cannot update repository - stopping beholder", slog.Any("error", lastError), logKeyRepositoryName, w.configuration.Name, logKeyCheckoutLocation, w.checkoutLocation, logKeyRepositoryURL, w.configuration.URL)
	w.broadcast(time.Now(), time.Now(), lastError)
	w.mutex.Unlock()
	w.Active = false
}

// broadcast issues the provided info via the heartbeat channel
func (w *Beholder) broadcast(last time.Time, next time.Time, error error) {
	w.heartbeat <- *&UpdateEvent{
		LastUpdate:     last,
		NextUpdate:     next,
		LastError:      error,
		RepositoryName: w.configuration.Name,
	}
}

// File retrieves the requested path from the managed repository
// File will ensure no file can be read if the repository is being updated
func (w *Beholder) File(filepath string) ([]byte, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	target := path.Join(w.checkoutLocation, path.Clean(filepath))
	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		return nil, fmt.Errorf("'%s' is not accessible or is a directory", target)
	}

	content, err := os.ReadFile(target)
	if err != nil || info.IsDir() {
		return nil, fmt.Errorf("an unexpected error occured while reading '%s' : %w", target, err)
	}

	return content, nil
}
