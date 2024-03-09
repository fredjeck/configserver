package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProblemDetail struct {
	ProblemType string `json:"type"`
	Title       string `json:"title"`
	Detail      string `json:"detail"`
	Instance    string `json:"instance"`
	Status      int    `json:"status"`
}

func HttpInternalServerError(w http.ResponseWriter, detail string, params ...interface{}) {
	writeStatus(w, http.StatusInternalServerError, "Internal Server Error", detail, params...)
}

func HttpUnauthorized(w http.ResponseWriter, detail string, params ...interface{}) {
	writeStatus(w, http.StatusUnauthorized, "Forbidden", detail, params...)
}

func HttpUnsupportedMediaType(w http.ResponseWriter, detail string, params ...interface{}) {
	writeStatus(w, http.StatusUnsupportedMediaType, "Unsupported content type", detail, params...)
}

func writeStatus(w http.ResponseWriter, code int, title string, detail string, params ...interface{}) {
	problem := &ProblemDetail{
		Status: code,
		Title:  title,
		Detail: fmt.Sprintf(detail, params...),
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(problem.Status)
	jsn, _ := json.Marshal(problem)
	_, _ = w.Write(jsn[0:len(jsn):len(jsn)])
}

func Ok(w http.ResponseWriter, content []byte) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content[0:len(content):len(content)])
}
