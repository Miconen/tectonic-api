package utils

import (
	"log/slog"
	"os"
	"strings"
)

var log = NewLogger()

func NewLogger() *slog.Logger {
	level_env := strings.ToLower(os.Getenv("LOG_LEVEL"))
	level := slog.LevelInfo

	switch level_env {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	return slog.New(slog.NewTextHandler(os.Stderr,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		},
	))
}
