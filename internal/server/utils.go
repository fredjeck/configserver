package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

// ContextRequestId key holding the request unique id in request contexts
const ContextRequestId = "http.request_id"

func die(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func dief(w http.ResponseWriter, statusCode int, message string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(message, args...), statusCode)
}

func dieErr(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {

	ctx := r.Context()
	reqIDRaw := ctx.Value(ContextRequestId) // reqIDRaw at this point is of type 'interface{}'
	reqId, ok := reqIDRaw.(string)
	if !ok {
		reqId = "not available"
	}

	slog.Error(message, "error", err, ContextRequestId, reqId)
	http.Error(w, message, statusCode)
}
