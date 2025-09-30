package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		slog.Info("Received new request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrappedWriter.statusCode,
			"duration", duration.String(),
			"duration_ms", duration.Milliseconds(),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
