package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// Adapted from https://blog.questionable.services/article/guide-logging-middleware-go/

const HttpRequestId = "http.request.id"
const HttpRequestStatus = "http.request.status"
const HttpRequestMethod = "http.request.method"
const HttpRequestPath = "http.request.path"
const HttpRequestDuration = "http.request.duration"

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

// RequestLoggingMiddleware logs the incoming HTTP request & its duration.
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
			ctx = context.WithValue(ctx, HttpRequestId, id.String())

			r = r.WithContext(ctx)

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			elapsed := time.Since(start)
			message := fmt.Sprintf("%s '%s' responded %d in %s", r.Method, r.URL.EscapedPath(), wrapped.status, elapsed)
			slog.Info(message,
				slog.Int(HttpRequestStatus, wrapped.status),
				slog.String(HttpRequestMethod, r.Method),
				slog.String(HttpRequestPath, r.URL.EscapedPath()),
				slog.Duration(HttpRequestDuration, elapsed),
				slog.String(HttpRequestId, id.String()),
			)
		}
		return http.HandlerFunc(fn)
	}
}
