package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func capture(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	old := logger
	logger = slog.New(slog.NewJSONHandler(&buf, nil))
	t.Cleanup(func() { logger = old })
	return &buf
}

func lines(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, l := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(l), &m); err != nil {
			t.Fatalf("bad log line %q: %v", l, err)
		}
		out = append(out, m)
	}
	return out
}

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(RequestLogger)
	r.Get("/guilds/{guild_id}/users/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		From(r.Context()).Info("handler line")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":4004,"message":"User not found"}`))
	})
	return r
}

func TestRequestIDPropagates(t *testing.T) {
	buf := capture(t)
	req := httptest.NewRequest(http.MethodGet, "/guilds/1/users/2", nil)
	req.Header.Set(RequestIDHeader, "1234567890")
	req.Header.Set(ActorHeader, "555")
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)

	if got := res.Header().Get(RequestIDHeader); got != "1234567890" {
		t.Errorf("echoed id %q", got)
	}
	ls := lines(t, buf)
	if len(ls) != 2 {
		t.Fatalf("expected handler + access line, got %d", len(ls))
	}
	for _, l := range ls {
		if l["request_id"] != "1234567890" || l["actor_id"] != "555" {
			t.Errorf("missing ids: %v", l)
		}
	}
	access := ls[1]
	if access["route"] != "/guilds/{guild_id}/users/{user_id}" || access["status"] != float64(404) ||
		access["error_code"] != float64(4004) || access["error"] != "User not found" || access["level"] != "WARN" {
		t.Errorf("access line: %v", access)
	}
}

func TestRequestIDGeneratedWhenMissingOrInvalid(t *testing.T) {
	for _, in := range []string{"", "bad id with spaces\n"} {
		capture(t)
		req := httptest.NewRequest(http.MethodGet, "/guilds/1/users/2", nil)
		if in != "" {
			req.Header.Set(RequestIDHeader, in)
		}
		res := httptest.NewRecorder()
		router().ServeHTTP(res, req)
		if got := res.Header().Get(RequestIDHeader); !validID.MatchString(got) || got == in {
			t.Errorf("input %q: got id %q", in, got)
		}
	}
}
