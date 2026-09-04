package http

import "net/http"

const allowedDevelopmentOrigin = "http://localhost:5173"

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == allowedDevelopmentOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type",
			)
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			if origin != allowedDevelopmentOrigin {
				http.Error(w, "CORS origin not allowed", http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
