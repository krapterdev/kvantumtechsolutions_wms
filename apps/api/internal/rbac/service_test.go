package rbac

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestNewService(t *testing.T) {
	service := NewService(nil)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if service.queries != nil {
		t.Fatal("NewService() should initialize with nil queries when nil is provided")
	}
}

func TestCreateRoleWithoutQueries(t *testing.T) {
	service := NewService(nil)

	organizationID := pgtype.UUID{
		Valid: true,
	}

	_, err := service.CreateRole(
		context.Background(),
		organizationID,
		"Administrator",
		"admin",
	)

	if err == nil {
		t.Fatal("CreateRole() expected error when queries are nil")
	}
}

func TestCreateRoleRequiresOrganizationID(t *testing.T) {
	service := NewService(nil)

	_, err := service.CreateRole(
		context.Background(),
		pgtype.UUID{},
		"Administrator",
		"admin",
	)

	if err == nil {
		t.Fatal("CreateRole() expected error for invalid organization ID")
	}
}

func TestCreateRoleRequiresName(t *testing.T) {
	service := NewService(nil)

	organizationID := pgtype.UUID{
		Valid: true,
	}

	_, err := service.CreateRole(
		context.Background(),
		organizationID,
		"",
		"admin",
	)

	if err == nil {
		t.Fatal("CreateRole() expected error for empty role name")
	}
}

func TestCreateRoleRequiresCode(t *testing.T) {
	service := NewService(nil)

	organizationID := pgtype.UUID{
		Valid: true,
	}

	_, err := service.CreateRole(
		context.Background(),
		organizationID,
		"Administrator",
		"",
	)

	if err == nil {
		t.Fatal("CreateRole() expected error for empty role code")
	}
}

func TestAssignRoleToUserRequiresUserID(t *testing.T) {
	service := NewService(nil)

	err := service.AssignRoleToUser(
		context.Background(),
		pgtype.UUID{},
		pgtype.UUID{Valid: true},
	)

	if err == nil {
		t.Fatal("AssignRoleToUser() expected error for invalid user ID")
	}
}

func TestAssignRoleToUserRequiresRoleID(t *testing.T) {
	service := NewService(nil)

	err := service.AssignRoleToUser(
		context.Background(),
		pgtype.UUID{Valid: true},
		pgtype.UUID{},
	)

	if err == nil {
		t.Fatal("AssignRoleToUser() expected error for invalid role ID")
	}
}

func TestUserHasPermissionRequiresUserID(t *testing.T) {
	service := NewService(nil)

	hasPermission, err := service.UserHasPermission(
		context.Background(),
		pgtype.UUID{},
		"users.read",
	)

	if err == nil {
		t.Fatal("UserHasPermission() expected error for invalid user ID")
	}

	if hasPermission {
		t.Fatal("UserHasPermission() should return false on error")
	}
}

func TestUserHasPermissionRequiresPermissionCode(t *testing.T) {
	service := NewService(nil)

	hasPermission, err := service.UserHasPermission(
		context.Background(),
		pgtype.UUID{Valid: true},
		"",
	)

	if err == nil {
		t.Fatal("UserHasPermission() expected error for empty permission code")
	}

	if hasPermission {
		t.Fatal("UserHasPermission() should return false on error")
	}
}
