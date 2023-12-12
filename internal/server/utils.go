package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func die(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func dief(w http.ResponseWriter, statusCode int, message string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(message, args...), statusCode)
}

func dieErr(w http.ResponseWriter, statusCode int, message string, err error) {
	slog.Error(message, "error", err)
	http.Error(w, message, statusCode)
}
