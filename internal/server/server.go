package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/savrin-sharif/go-rest-template/internal/config"
)

// Server wraps the HTTP server and lifecycle hooks.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// New constructs a configured HTTP server.
func New(cfg config.Config, logger *slog.Logger) *Server {
	router := newRouter(cfg, logger)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{httpServer: httpServer, logger: logger}
}

// Start runs the HTTP server.
func (s *Server) Start() error {
	s.logger.Info("http server starting", slog.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("http server shutting down")
	return s.httpServer.Shutdown(ctx)
}
