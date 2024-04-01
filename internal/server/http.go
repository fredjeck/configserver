package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ctxRequestID struct{}
type ctxClientID struct{}

// ProblemDetail is a RFC9457 compliant error detail used by the server to return errors.
type ProblemDetail struct {
	ProblemType string `json:"type"`
	Title       string `json:"title"`
	Detail      string `json:"detail"`
	Instance    string `json:"instance"`
	Status      int    `json:"status"`
}

// HTTPInternalServerError returns an HTTP 500 error along a RFC9457 compliant error detail
func HTTPInternalServerError(w http.ResponseWriter, r *http.Request, detail string, params ...interface{}) {
	writeStatus(w, r, http.StatusInternalServerError, "Internal Server Error", detail, params...)
}

// HTTPUnauthorized returns an HTTP 401 error along a RFC9457 compliant error detail
func HTTPUnauthorized(w http.ResponseWriter, r *http.Request, detail string, params ...interface{}) {
	writeStatus(w, r, http.StatusUnauthorized, "Forbidden", detail, params...)
}

// HTTPUnsupportedMediaType returns an HTTP 415 error along a RFC9457 compliant error detail
func HTTPUnsupportedMediaType(w http.ResponseWriter, r *http.Request, detail string, params ...interface{}) {
	writeStatus(w, r, http.StatusUnsupportedMediaType, "Unsupported content type", detail, params...)
}

// HTTPNotFound returns an HTTP 404 error along a RFC9457 compliant error detail
func HTTPNotFound(w http.ResponseWriter, r *http.Request, detail string, params ...interface{}) {
	writeStatus(w, r, http.StatusUnsupportedMediaType, "Not found", detail, params...)
}

func writeStatus(w http.ResponseWriter, r *http.Request, code int, title string, detail string, params ...interface{}) {

	strDetail := fmt.Sprintf(detail, params...)
	if code > 300 {
		requestID := r.Context().Value(ctxRequestID{}).(string)
		slog.Warn(strDetail, HTTPRequestStatus, code, HTTPRequestID, requestID)
	}
	problem := &ProblemDetail{
		Status: code,
		Title:  title,
		Detail: strDetail,
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(problem.Status)
	jsn, _ := json.Marshal(problem)
	_, _ = w.Write(jsn[0:len(jsn):len(jsn)])
}

// Ok returns an HTTP 201 response along the provided content
func Ok(w http.ResponseWriter, content []byte, mimetype string) {
	w.Header().Add("Content-Type", mimetype)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content[0:len(content):len(content)])
}
