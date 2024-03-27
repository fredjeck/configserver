package server

import (
	"encoding/json"
	"net/http"

	"github.com/fredjeck/configserver/internal/repository"
)

// Handles the clients file tokenization requests
func handleStatistics(mgr *repository.Manager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsn, _ := json.Marshal(mgr.Statistics())
		Ok(w, jsn, "application/json;charset=utf-8")
	}
}
