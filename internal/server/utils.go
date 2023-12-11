package server

import (
	"log/slog"
	"net/http"
)

func die(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func dieErr(w http.ResponseWriter, statusCode int, message string, err error) {
	slog.Error(message, "error", err)
	http.Error(w, message, statusCode)
}
