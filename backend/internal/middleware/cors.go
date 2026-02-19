package middleware

import (
	"net/http"
	"os"
)

const (
	defaultCORSOrigin = "http://localhost:3000"
	corsOriginEnv     = "CORS_ORIGIN"
)

// CORS returns middleware that sets Cross-Origin Resource Sharing headers.
// The allowed origin is read from the CORS_ORIGIN environment variable,
// defaulting to http://localhost:3000.
func CORS() func(http.Handler) http.Handler {
	origin := os.Getenv(corsOriginEnv)
	if origin == "" {
		origin = defaultCORSOrigin
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
