package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

// Adapted from https://blog.questionable.services/article/guide-logging-middleware-go/

// HTTPRequestID represents the logging key for the http request Id
const HTTPRequestID = "http.request.id"

// HTTPRequestStatus represents the logging key for the http request Status
const HTTPRequestStatus = "http.request.status"

// HTTPRequestMethod represents the logging key for the http request Method
const HTTPRequestMethod = "http.request.method"

// HTTPRequestPath represents the logging key for the http request Path
const HTTPRequestPath = "http.request.path"

// HTTPRequestDuration represents the logging key for the http request Duration
const HTTPRequestDuration = "http.request.duration"

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// requestLogger logs the incoming HTTP request & its duration.
func requestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					slog.Error("internal server error", slog.String("stacktrace", string(debug.Stack())))
				}
			}()

			start := time.Now()

			ctx := r.Context()
			id, _ := uuid.NewV7()
			ctx = context.WithValue(ctx, ctxRequestID{}, id.String())

			r = r.WithContext(ctx)

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			elapsed := time.Since(start)
			message := fmt.Sprintf("%s '%s' responded %d in %s", r.Method, r.URL.EscapedPath(), wrapped.status, elapsed)
			slog.Info(message,
				slog.Int(HTTPRequestStatus, wrapped.status),
				slog.String(HTTPRequestMethod, r.Method),
				slog.String(HTTPRequestPath, r.URL.EscapedPath()),
				slog.Duration(HTTPRequestDuration, elapsed),
				slog.String(HTTPRequestID, id.String()),
			)
		}
		return http.HandlerFunc(fn)
	}
}
