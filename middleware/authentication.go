package middleware

import (
	"net/http"
	"os"
	"tectonic-api/utils"
)

var log = utils.NewLogger()

func Authentication(next http.Handler) http.Handler {
	log.Debug("Adding authentication handler")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rlog := log.With(
			"method", r.Method,
			"url", r.URL,
		)

		apiKey := r.Header.Get("Authorization")

		rlog.Debug("Validating API key")
		validApiKey := os.Getenv("API_KEY")
		if apiKey != validApiKey {
			rlog.Warn("Authentication key is invalid", "key", apiKey)
			jw := utils.NewJsonWriter(w, r, http.StatusUnauthorized)
			jw.WriteResponse(http.NoBody)
			return
		}

		rlog.Debug("API key is valid")
		next.ServeHTTP(w, r)
	})
}
