package middleware

import (
	"context"
	"fmt"
	"github.com/fredjeck/configserver/internal/server"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// Adapted from https://blog.questionable.services/article/guide-logging-middleware-go/

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
func RequestLoggingMiddleware() func(http.Handler) http.Handler {
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
			id := uuid.New()
			ctx = context.WithValue(ctx, server.ContextRequestId, id.String())

			r = r.WithContext(ctx)

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			elapsed := time.Since(start)
			message := fmt.Sprintf("%s '%s' responded %d in %s", r.Method, r.URL.EscapedPath(), wrapped.status, elapsed)
			slog.Info(message,
				slog.Int("http.status", wrapped.status),
				slog.String("http.method", r.Method),
				slog.String("http.path", r.URL.EscapedPath()),
				slog.Duration("http.duration", elapsed),
				slog.String(server.ContextRequestId, id.String()),
			)
		}
		return http.HandlerFunc(fn)
	}
}
