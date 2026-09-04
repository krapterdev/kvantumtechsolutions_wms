package auth

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.sessionManager == nil || h.service == nil {
			http.Error(w, "authentication service unavailable", http.StatusInternalServerError)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		session, err := h.sessionManager.GetSessionByToken(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		user, err := h.service.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		if !user.IsActive {
			http.Error(w, "user is inactive", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}
