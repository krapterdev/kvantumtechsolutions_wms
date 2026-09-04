package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/auth"
	"github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/database"
	db "github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/database/generated"
	"github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/rbac"
)

const testOrganizationID = "fbfc357f-6973-46b4-9db5-bc731140f9cd"

func TestRouterCreateUserRequiresAuthentication(t *testing.T) {
	router := NewRouter(nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/users",
		nil,
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRouterCreateUserPermissionDenied(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	organizationID := pgtype.UUID{}
	if err := organizationID.Scan(testOrganizationID); err != nil {
		t.Fatalf("parse organization ID: %v", err)
	}

	userID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}

	roleID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}

	permissionID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}

	permissionCode := "router.test.create." + uuid.New().String()
	testEmail := "router-rbac-test-" + uuid.New().String() + "@example.com"

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cleanupCancel()

		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM sessions WHERE user_id = $1",
			userID,
		)

		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM user_roles WHERE user_id = $1",
			userID,
		)

		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM role_permissions WHERE role_id = $1",
			roleID,
		)

		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM users WHERE id = $1",
			userID,
		)

		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM roles WHERE id = $1",
			roleID,
		)

		_, _ = pool.Exec(
			cleanupCtx,
			"DELETE FROM permissions WHERE id = $1",
			permissionID,
		)
	})

	permission, err := queries.CreatePermission(
		ctx,
		db.CreatePermissionParams{
			ID:   permissionID,
			Code: permissionCode,
			Name: "Router Test Create",
			Description: pgtype.Text{
				String: "Temporary permission for router integration test",
				Valid:  true,
			},
		},
	)
	if err != nil {
		t.Fatalf("create test permission: %v", err)
	}

	role, err := queries.CreateRole(
		ctx,
		db.CreateRoleParams{
			ID:             roleID,
			OrganizationID: organizationID,
			Name:           "Router Test Role",
			Code:           "router-test-" + uuid.New().String()[:8],
		},
	)
	if err != nil {
		t.Fatalf("create test role: %v", err)
	}

	if err := queries.CreateRolePermission(
		ctx,
		db.CreateRolePermissionParams{
			RoleID:       role.ID,
			PermissionID: permission.ID,
		},
	); err != nil {
		t.Fatalf("create role permission: %v", err)
	}

	passwordHash, err := auth.HashPassword("TestPassword123!")
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}

	testUser, err := queries.CreateUser(
		ctx,
		db.CreateUserParams{
			ID:             userID,
			OrganizationID: organizationID,
			Email:          testEmail,
			PasswordHash:   passwordHash,
			FirstName:      "Router",
			LastName:       pgtype.Text{},
			IsActive:       true,
		},
	)
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}

	if err := queries.AssignRoleToUser(
		ctx,
		db.AssignRoleToUserParams{
			UserID: userID,
			RoleID: roleID,
		},
	); err != nil {
		t.Fatalf("assign test role: %v", err)
	}

	sessionManager := auth.NewSessionManager(queries)

	sessionToken, err := sessionManager.CreateSession(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("create test session: %v", err)
	}

	router := NewRouter(queries)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/users",
		nil,
	)

	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: sessionToken,
	})

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestRBACServiceIsAvailable(t *testing.T) {
	service := rbac.NewService(nil)

	if service == nil {
		t.Fatal("expected RBAC service to be created")
	}
}
