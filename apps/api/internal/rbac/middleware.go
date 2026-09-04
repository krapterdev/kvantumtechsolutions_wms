package rbac

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/auth"
)

func RequirePermission(
	service *Service,
	permissionCode string,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(
				w,
				"rbac service unavailable",
				http.StatusInternalServerError,
			)
			return
		}

		user, ok := auth.UserFromContext(r.Context())
		if !ok {
			http.Error(
				w,
				"authentication required",
				http.StatusUnauthorized,
			)
			return
		}

		userID := pgtype.UUID{
			Bytes: user.ID.Bytes,
			Valid: user.ID.Valid,
		}

		hasPermission, err := service.UserHasPermission(
			r.Context(),
			userID,
			permissionCode,
		)
		if err != nil {
			http.Error(
				w,
				"authorization check failed",
				http.StatusInternalServerError,
			)
			return
		}

		if !hasPermission {
			http.Error(
				w,
				"permission denied",
				http.StatusForbidden,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
