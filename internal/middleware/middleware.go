package middleware

import (
	"bitaksi-go-driver/internal/config"
	"net/http"
)

// APIKeyMiddleware validates the API key in the Authorization header
func APIKeyMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the Authorization header
			apiKey := r.Header.Get("Authorization")

			// Validate the API key
			if apiKey != cfg.Server.APIKey {
				http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
				return
			}

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
