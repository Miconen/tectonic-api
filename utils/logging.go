package utils

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
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

		if os.Getenv("RAILWAY_PROJECT_ID") != "" {
			log.Info(r.URL.Path,
				"header", r.Header,
				"method", r.Method,
				"status", recorder.statusCode,
				"path", r.URL,
				"time", start,
				"host", r.Host,
				"duration", duration.String(),
			)
		} else {
			log.Info(r.URL.Path,
				"method", r.Method,
				"status", recorder.statusCode,
				"time", start,
				"duration", duration.String(),
			)
		}

	})
}
