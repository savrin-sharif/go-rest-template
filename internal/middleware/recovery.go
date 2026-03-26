package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/savrin-sharif/go-rest-template/pkg/httputil"
)

// Recover intercepts panics, logs them, and returns a 500 response.
func Recover(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, rec any) {
		logger.Error("panic recovered",
			slog.Any("error", rec),
			slog.String("request_id", c.GetHeader("X-Request-ID")),
			slog.String("path", c.Request.URL.Path),
			slog.String("method", c.Request.Method),
			slog.String("stacktrace", string(debug.Stack())),
		)

		httputil.WriteJSON(c.Writer, 500, map[string]string{
			"error": "internal server error",
		})
		c.Abort()
	})
}
