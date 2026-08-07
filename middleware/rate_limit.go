package middleware

import (
	"encoding/json"
	"net/http"

	"tectonic-api/models"

	"golang.org/x/time/rate"
)

// Rate limiting
func RateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(120), 10)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			apiErr := models.NewTectonicError(models.ERROR_API_RATE_LIMITED)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(apiErr.GetStatus())
			json.NewEncoder(w).Encode(apiErr)
			return
		}
		next.ServeHTTP(w, r)
	})
}
