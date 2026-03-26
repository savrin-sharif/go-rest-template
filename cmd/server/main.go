package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/savrin-sharif/go-rest-template/internal/config"
	"github.com/savrin-sharif/go-rest-template/internal/server"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		slog.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	logger := config.NewLogger(cfg.Log.Level, cfg.Log.AddSource)

	srv := server.New(cfg, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("server stopped unexpectedly", slog.Any("error", err))
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	// Allow in-flight logs to flush.
	time.Sleep(200 * time.Millisecond)
	logger.Info("server exited cleanly")
}
