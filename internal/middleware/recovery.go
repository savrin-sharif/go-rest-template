package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	chmiddleware "github.com/go-chi/chi/v5/middleware"
)

// Recover intercepts panics, logs them, and returns a 500 response.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					reqID := chmiddleware.GetReqID(r.Context())
					logger.Error("panic recovered",
						slog.Any("error", rec),
						slog.String("request_id", reqID),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.String("stacktrace", string(debug.Stack())),
					)

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
