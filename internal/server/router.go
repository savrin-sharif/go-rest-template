package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/savrin-sharif/go-rest-template/internal/config"
	"github.com/savrin-sharif/go-rest-template/internal/handler"
	appmw "github.com/savrin-sharif/go-rest-template/internal/middleware"
)

func newRouter(cfg config.Config, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(chmiddleware.RequestID)
	r.Use(chmiddleware.RealIP)
	r.Use(appmw.Recover(logger))
	r.Use(appmw.Logger(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	healthHandler := handler.NewHealthHandler(cfg.AppName, logger)
	r.Get("/", healthHandler.Welcome)
	r.Get("/health", healthHandler.Health)

	return r
}
