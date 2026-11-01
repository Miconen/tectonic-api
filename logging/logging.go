package logging

import (
	"log/slog"
	"os"
	"strings"

	"tectonic-api/config"
)

var (
	logger   *slog.Logger
	settings *LoggerConfig
)

// Get returns the global logger. Before Init (in tests, say) it falls back to
// slog's default instead of exiting.
func Get() *slog.Logger {
	if logger == nil {
		return slog.Default()
	}
	return logger
}

type LoggerConfig struct {
	level      string
	jsonOutput bool
}

func Init(cfg *config.Config) {
	settings = &LoggerConfig{
		level:      strings.ToLower(cfg.LogLevel),
		jsonOutput: false,
	}

	if cfg.RailwayProjectID != "" {
		settings.jsonOutput = true
	}

	level := slog.LevelInfo

	switch settings.level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("15:04:05"))
			}
			return a
		},
		Level: level,
	}

	if settings.jsonOutput {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
		return
	}

	logger = slog.New(slog.NewTextHandler(os.Stderr, opts))
}
