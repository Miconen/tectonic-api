package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
)

// Rate limiting
func RateLimit(limit int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(limit), 10)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
