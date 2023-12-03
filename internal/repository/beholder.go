package repository

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"log/slog"
	"os"
	"path"
	"sync"
	"time"
)

// Beholder are responsible for maintaining local copies of git repositories up to date based on the provided configuration
// As they are running in background they make use of a heartbeat channel towards their initiator to communicate about
// repositories update event
type Beholder struct {
	*Configuration                  // The current configured repository
	Active         bool             // true if the repository is actively monitored and can be used
	heartbeat      chan UpdateEvent // Uplink to the beholder's initiator
	mutex          *sync.RWMutex    // Used to ensure no read operation is allowed while the repository is being updated
}

// NewBeholder initiates a new beholder for the provided configuration
// a call to Watch() is mandatory to start the beholder process
func NewBeholder(configuration *Configuration, heartbeat chan UpdateEvent) *Beholder {
	return &Beholder{configuration, true, heartbeat, &sync.RWMutex{}}
}

// Watch initiates the creation of a local copy of the configured repository and will periodically update the repository
// to its latest state
func (w *Beholder) Watch() {
	go w.watchInternal()
}

// see Watch
func (w *Beholder) watchInternal() {
	var lastError error = nil

	for {
		w.mutex.Lock()
		last := time.Now()

		slog.Info("pulling repository", LogKeyRepositoryName, w.Name)

		if err := os.MkdirAll(w.CheckoutLocation, os.ModePerm); err != nil {
			lastError = fmt.Errorf("cannot create path '%s' to checkout '%s': %w", w.CheckoutLocation, w.Name, err)
			break
		}

		workspace, err := git.PlainOpen(w.CheckoutLocation)
		if err != nil {
			slog.Info("no local copy found, creating a fresh clone", LogKeyRepositoryName, w.Name)
			workspace, err = git.PlainClone(w.CheckoutLocation, false, &git.CloneOptions{
				URL:        w.Url,
				Progress:   os.Stdout,
				RemoteName: w.Branch,
			})

			if err != nil {
				lastError = fmt.Errorf("could not clone '%s' to '%s' : %w", w.Url, w.CheckoutLocation, err)
				break
			}
		}

		tree, err := workspace.Worktree()
		if err != nil {
			lastError = fmt.Errorf("'%s' : unable to open local copy : %s", w.CheckoutLocation, err)
			break
		}

		err = tree.Pull(&git.PullOptions{
			RemoteName: w.Branch,
			Force:      true,
		})

		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			lastError = fmt.Errorf("'%s' : unable to pull latest changes : %s", w.Name, err)
			break
		}

		nextRefresh := time.Duration(w.RefreshIntervalSeconds) * time.Second
		slog.Info(fmt.Sprintf("'%s' next pull will occur @ %s", w.Name, time.Now().Add(nextRefresh)), LogKeyRepositoryName, w.Name)
		w.broadcast(last, time.Now().Add(nextRefresh), nil)
		w.mutex.Unlock()
		time.Sleep(nextRefresh)
	}

	// If we are here something bad happend
	slog.Error("Cannot update repository - stopping beholder", slog.Any("error", lastError))
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
		RepositoryName: w.Name,
	}
}

// File retrieves the requested path from the managed repository
// File will ensure no file can be read if the repository is being updated
func (w *Beholder) File(filepath string) ([]byte, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	target := path.Join(w.CheckoutLocation, path.Clean(filepath))
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
