package server

import (
	"encoding/json"
	"github.com/fredjeck/configserver/pkg/repo"
	"net/http"
)

type StatisticsResponse struct {
	repo.RepositoryStatistics
	Name string `json:"name"`
}

func (server *ConfigServer) statistics(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	w.Header().Add("Content-Type", "application/json")

	stats := make([]*StatisticsResponse, 0)

	for _, s := range server.repositories.Repositories {
		stats = append(stats, &StatisticsResponse{
			RepositoryStatistics: *s.Statistics,
			Name:                 s.Configuration.Name,
		})
	}

	values, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	server.writeResponse(http.StatusOK, values, w)
}
