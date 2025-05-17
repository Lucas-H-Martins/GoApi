package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	// This is just an example token. In a real application, you would use a proper token management system
	VALID_TOKEN = "your-secret-token-123"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "No authorization header provided",
			})
			return
		}

		// Check if the auth header starts with "Bearer "
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid authorization header format",
			})
			return
		}

		// Verify the token
		token := bearerToken[1]
		if token != VALID_TOKEN {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid token",
			})
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
} 