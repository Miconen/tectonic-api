package middleware

import (
	"net/http"
	"os"
	"tectonic-api/utils"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the API key from the URL query parameters
		apiKey := r.URL.Query().Get("api_key")

		// Validate the API key (you may replace this with your own validation logic)
		validApiKey := os.Getenv("API_KEY")
		if apiKey != validApiKey {
			utils.JsonWriter(http.NoBody).IntoHTTP(http.StatusUnauthorized)(w, r)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
