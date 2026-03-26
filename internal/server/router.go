package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-rest-template/internal/config"
	"github.com/savrin-sharif/go-rest-template/internal/handler"
)

func newRouter(cfg config.Config, logger *slog.Logger) http.Handler {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposeHeaders:    []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           5 * time.Minute,
	}))

	healthHandler := handler.NewHealthHandler(cfg.AppName, logger)
	r.GET("/", gin.WrapF(healthHandler.Welcome))
	r.GET("/health", gin.WrapF(healthHandler.Health))

	return r
}
