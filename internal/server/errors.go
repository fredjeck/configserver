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

func HttpNotAuthorized(w http.ResponseWriter, detail string, params ...interface{}) {
	problem := &ProblemDetail{
		Status: http.StatusUnauthorized,
		Title:  "Unauthorized",
		Detail: fmt.Sprintf(detail, params...),
	}
	writeStatus(w, problem)
}

func HttpBadRequest(w http.ResponseWriter, detail string, params ...interface{}) {
	problem := &ProblemDetail{
		Status: http.StatusBadRequest,
		Title:  "Bad Request",
		Detail: fmt.Sprintf(detail, params...),
	}
	writeStatus(w, problem)
}

func HttpNotFound(w http.ResponseWriter, detail string, params ...interface{}) {
	problem := &ProblemDetail{
		Status: http.StatusNotFound,
		Title:  "Not Found",
		Detail: fmt.Sprintf(detail, params...),
	}
	writeStatus(w, problem)
}

func HttpInternalServerError(w http.ResponseWriter, detail string, params ...interface{}) {
	problem := &ProblemDetail{
		Status: http.StatusInternalServerError,
		Title:  "Not Found",
		Detail: fmt.Sprintf(detail, params...),
	}
	writeStatus(w, problem)
}

func writeStatus(w http.ResponseWriter, detail *ProblemDetail) {
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(detail.Status)
	jsn, _ := json.Marshal(detail)
	_, _ = fmt.Fprint(w, jsn)
}
