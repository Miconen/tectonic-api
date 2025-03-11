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

	if os.Getenv("RAILWAY_PROJECT_ID") != "" {
		return slog.New(slog.NewJSONHandler(os.Stderr,
			&slog.HandlerOptions{
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						t := a.Value.Time()
						a.Value = slog.StringValue(t.Format("15:04:05"))
					}
					return a
				},
				AddSource: true,
				Level:     level,
			},
		))
	} else {
		return slog.New(slog.NewTextHandler(os.Stderr,
			&slog.HandlerOptions{
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						t := a.Value.Time()
						a.Value = slog.StringValue(t.Format("15:04:05"))
					}
					return a
				},
				AddSource: true,
				Level:     level,
			},
		))

	}

}
