package config

import (
	"log/slog"
	"os"
)

// NewLogger constructs a slog Logger with sensible defaults.
func NewLogger(levelStr string, addSource bool) *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	})

	return slog.New(handler)
}
