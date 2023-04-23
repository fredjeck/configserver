package server

import (
	"encoding/json"
	"net/http"
)

func (server *ConfigServer) listRepositories(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	w.Header().Add("Content-Type", "application/json")
	var repos []string
	for _, v := range server.configuration.Repositories {
		repos = append(repos, v.Name)
	}
	values, err := json.Marshal(repos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	server.writeResponse(http.StatusOK, values, w)
}
