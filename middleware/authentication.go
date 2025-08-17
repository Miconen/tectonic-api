package middleware

import (
	"net/http"
	"tectonic-api/config"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

func Authentication(cfg *config.Config) mux.MiddlewareFunc {
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
				jw := utils.NewJsonWriter(w, r, http.StatusUnauthorized)
				jw.WriteError(models.ERROR_INVALID_TOKEN)
				return
			}

			rlog.Debug("API key is valid")
			next.ServeHTTP(w, r)
		})
	}
}
