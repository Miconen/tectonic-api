package middleware

import (
	"net/http"
	"os"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"
)

func Authentication(next http.Handler) http.Handler {
	logging.Get().Debug("Adding authentication handler")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rlog := logging.Get().With(
			"method", r.Method,
			"url", r.URL,
		)

		apiKey := r.Header.Get("Authorization")

		rlog.Debug("Validating API key")
		validApiKey := os.Getenv("API_KEY")
		if apiKey != validApiKey {
			rlog.Warn("Authentication key is invalid")
			jw := utils.NewJsonWriter(w, r, http.StatusUnauthorized)
			jw.WriteError(models.ERROR_INVALID_TOKEN)
			return
		}

		rlog.Debug("API key is valid")
		next.ServeHTTP(w, r)
	})
}
