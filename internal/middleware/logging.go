package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger adds structured request logging with request ID correlation.
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if q := c.Request.URL.RawQuery; q != "" {
			path = fmt.Sprintf("%s?%s", path, q)
		}

		c.Next()

		logger.With(
			slog.String("request_id", c.GetHeader("X-Request-ID")),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.Int("bytes", c.Writer.Size()),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		).Info("request completed")
	}
}
