package server

import (
	"fmt"
	"go.uber.org/zap"
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
func RequestLoggingMiddleware(logger zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("Internal server error", zap.ByteString("trace", debug.Stack()))

				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			elapsed := time.Since(start)
			message := fmt.Sprintf("%s %s %d %s", r.Method, r.URL.EscapedPath(), wrapped.status, elapsed)
			logger.Info(message,
				zap.Int("http.status", wrapped.status),
				zap.String("http.method", r.Method),
				zap.String("http.path", r.URL.EscapedPath()),
				zap.Duration("http.duration", elapsed),
			)
		}
		return http.HandlerFunc(fn)
	}
}
