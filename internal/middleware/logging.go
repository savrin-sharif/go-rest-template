package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chmiddleware "github.com/go-chi/chi/v5/middleware"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// Logger adds structured request logging with request ID correlation.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := chmiddleware.GetReqID(r.Context())

			rw := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r)

			logger.With(
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", statusFromResponse(rw)),
				slog.Int("bytes", rw.bytes),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			).Info("request completed")
		})
	}
}

func statusFromResponse(rw *responseWriter) int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}
