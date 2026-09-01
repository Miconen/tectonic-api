package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"tectonic-api/config"
	"tectonic-api/logging"
	"tectonic-api/models"
)

func isReadOnlyRequest(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodHead
}

func writeAuthenticationError(w http.ResponseWriter, code models.APIV1ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code.Status())
	_ = json.NewEncoder(w).Encode(models.NewTectonicError(code))
}

func Authentication(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logging.Get().Debug("Adding authentication handler")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for docs and OpenAPI spec
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			rlog := logging.Get().With("method", r.Method, "url", r.URL)

			token := r.Header.Get("Authorization")

			rlog.Debug("Validating API key")
			if token == cfg.APIKey {
				rlog.Debug("Full-access API key is valid")
				next.ServeHTTP(w, r)
				return
			}

			if token == cfg.APIReadKey {
				if !isReadOnlyRequest(r) {
					rlog.Warn("Read-only API key attempted a write operation")
					writeAuthenticationError(w, models.ERROR_INSUFFICIENT_SCOPE)
					return
				}

				rlog.Debug("Read-only API key is valid")
				next.ServeHTTP(w, r)
				return
			}

			rlog.Warn("Authentication key is invalid")
			writeAuthenticationError(w, models.ERROR_INVALID_TOKEN)
		})
	}
}
