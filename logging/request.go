package logging

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
)

// Headers that tie a request to the caller's own logs (#99).
const (
	// RequestIDHeader carries the caller's correlation ID (the bot uses the
	// Discord interaction ID). One is generated when it's missing, and it's
	// always echoed back.
	RequestIDHeader = "X-Request-ID"
	// ActorHeader is the Discord user who triggered the request, so one
	// member's actions can be followed across requests.
	ActorHeader = "X-Actor-ID"
)

type ctxKey struct{}

// From returns the request's logger, which tags every line with request_id
// (and actor_id when the caller sent one). Outside a request it returns the
// global logger.
func From(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return Get()
}

// RequestID returns the current request's ID, or "".
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

type requestIDKey struct{}

var validID = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,64}$`)

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// RequestLogger must be the first middleware. It sets up the request ID and
// logger, then writes one access log line per request with the route, status,
// duration and, for failures, the API error code and message.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if !validID.MatchString(id) {
			id = newID()
		}
		w.Header().Set(RequestIDHeader, id)

		attrs := []any{"request_id", id}
		if actor := r.Header.Get(ActorHeader); validID.MatchString(actor) {
			attrs = append(attrs, "actor_id", actor)
		}
		log := Get().With(attrs...)
		ctx := context.WithValue(r.Context(), ctxKey{}, log)
		ctx = context.WithValue(ctx, requestIDKey{}, id)

		rec := &recorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r.WithContext(ctx))

		accessLog(log, r, rec, time.Since(start))
	})
}

func accessLog(log *slog.Logger, r *http.Request, rec *recorder, d time.Duration) {
	route := r.URL.Path
	if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
		route = rc.RoutePattern()
	}
	attrs := []any{
		"method", r.Method,
		"route", route,
		"path", r.URL.Path,
		"status", rec.status,
		"duration_ms", d.Milliseconds(),
		"user_agent", r.UserAgent(),
	}

	level := slog.LevelInfo
	if rec.status >= 400 {
		level = slog.LevelWarn
		var body struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Detail  string `json:"detail"`
		}
		if json.Unmarshal(rec.errBody.Bytes(), &body) == nil {
			msg := body.Message
			if msg == "" {
				msg = body.Detail
			}
			attrs = append(attrs, "error_code", body.Code, "error", msg)
		}
	}
	if rec.status >= 500 {
		level = slog.LevelError
	}
	log.Log(context.Background(), level, r.Method+" "+route+" "+http.StatusText(rec.status), attrs...)
}

// recorder captures the status code, and the first 2 KB of error bodies so
// the access log can say why a request failed.
type recorder struct {
	http.ResponseWriter
	status  int
	errBody bytes.Buffer
}

func (r *recorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *recorder) Write(b []byte) (int, error) {
	if r.status >= 400 && r.errBody.Len() < 2048 {
		r.errBody.Write(b[:min(len(b), 2048-r.errBody.Len())])
	}
	return r.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the real writer.
func (r *recorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }
