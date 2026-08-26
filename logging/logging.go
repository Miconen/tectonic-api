package logging

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"tectonic-api/config"
)

var (
	logger   *slog.Logger
	settings *LoggerConfig
)

func Get() *slog.Logger {
	if logger == nil {
		fmt.Fprintf(os.Stderr, "Error loading logger")
		os.Exit(1)
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

// statusRecorder wraps the ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before calling the wrapped WriteHeader
func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		h.ServeHTTP(recorder, r)
		duration := time.Since(start)

		if settings.jsonOutput {
			Get().Info(
				r.Method+" "+r.URL.Path,
				"status", recorder.statusCode,
				"duration", duration.String(),
				"host", r.Host,
				"origin", r.Header.Get("Origin"),
				"user_agent", r.UserAgent(),
				"content_type", r.Header.Get("Content-Type"),
				"method", r.Method,
				"path", r.URL.Path,
			)
		} else {
			Get().Info(
				r.Method+" "+r.URL.Path,
				"method", r.Method,
				"status", recorder.statusCode,
				"time", start,
				"duration", duration.String(),
			)
		}
	})
}
