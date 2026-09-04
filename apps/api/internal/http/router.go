package http

import (
	"encoding/json"
	"net/http"

	"github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/auth"
	db "github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/database/generated"
)

func NewRouter(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	authService := auth.NewService(queries)
	sessionManager := auth.NewSessionManager(queries)
	authHandler := auth.NewHandler(authService, sessionManager)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","version":"v1"}`))
	})

	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/users", authHandler.CreateUser)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	mux.Handle("GET /api/v1/auth/me", authHandler.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := auth.UserFromContext(r.Context())
			if !ok {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			response := struct {
				User struct {
					ID             string `json:"id"`
					OrganizationID string `json:"organization_id"`
					Email          string `json:"email"`
					FirstName      string `json:"first_name"`
					LastName       string `json:"last_name,omitempty"`
					IsActive       bool   `json:"is_active"`
				} `json:"user"`
			}{}

			response.User.ID = user.ID.String()
			response.User.OrganizationID = user.OrganizationID.String()
			response.User.Email = user.Email
			response.User.FirstName = user.FirstName
			response.User.IsActive = user.IsActive

			if user.LastName.Valid {
				response.User.LastName = user.LastName.String
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_ = json.NewEncoder(w).Encode(response)
		}),
	))

	return mux
}
