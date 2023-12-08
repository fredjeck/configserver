// Package repository contains utilities abstracting git repositories configuration
// and source code pulling
package repository

import (
	"time"
)

const LogKeyRepositoryName = "repository_name"

// Repository is a handle on a git repository
type Repository struct {
	Configuration *Configuration // repository configuration
	Statistics    *Statistics    // repository access statistics
	Beholder      *Beholder      // beholder process managing the repository
}

// Configuration represents a repository configuration as retrieved from yaml configuration files
type Configuration struct {
	Name                   string   `yaml:"name"`
	Url                    string   `yaml:"url"`
	Branch                 string   `yaml:"branch"`
	RefreshIntervalSeconds int      `yaml:"refreshIntervalSeconds"`
	CheckoutLocation       string   `yaml:"checkoutLocation"`
	Token                  string   `yaml:"token"`
	Clients                []string `yaml:"clients"`
}

// Statistics allows to maintain some stats about repository access
type Statistics struct {
	Hitcount   int64     `json:"hitCount"`
	LastUpdate time.Time `json:"lastUpdate"`
	NextUpdate time.Time `json:"nextUpdate"`
	LastError  error     `json:"lastError"`
}

// UpdateEvent as generated by beholders
type UpdateEvent struct {
	RepositoryName string
	LastUpdate     time.Time
	NextUpdate     time.Time
	LastError      error
	Active         bool
}

// IsClientAllowed verifies if the provided clientId is allowed to access the repository based on its configuration
func (repo *Repository) IsClientAllowed(clientId string) bool {
	for _, client := range repo.Configuration.Clients {
		if client == clientId {
			return true
		}
	}
	return false
}
