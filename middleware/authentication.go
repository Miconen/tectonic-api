package middleware

import (
	"encoding/json"
	"net/http"

	"tectonic-api/config"
	"tectonic-api/logging"
	"tectonic-api/models"
)

func Authentication(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logging.Get().Debug("Adding authentication handler")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rlog := logging.Get().With(
				"method", r.Method,
				"url", r.URL,
			)

			token := r.Header.Get("Authorization")

			rlog.Debug("Validating API key")
			if token != cfg.APIKey {
				rlog.Warn("Authentication key is invalid")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(models.NewTectonicError(models.ERROR_INVALID_TOKEN))
				return
			}

			rlog.Debug("API key is valid")
			next.ServeHTTP(w, r)
		})
	}
}
