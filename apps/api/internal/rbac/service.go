package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/database/generated"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrForbidden          = errors.New("permission denied")
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) CreateRole(
	ctx context.Context,
	organizationID pgtype.UUID,
	name string,
	code string,
) (db.Role, error) {
	if s.queries == nil {
		return db.Role{}, fmt.Errorf("rbac service queries are not initialized")
	}

	if !organizationID.Valid {
		return db.Role{}, fmt.Errorf("organization ID is required")
	}

	if name == "" {
		return db.Role{}, fmt.Errorf("role name is required")
	}

	if code == "" {
		return db.Role{}, fmt.Errorf("role code is required")
	}

	roleID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}

	return s.queries.CreateRole(ctx, db.CreateRoleParams{
		ID:             roleID,
		OrganizationID: organizationID,
		Name:           name,
		Code:           code,
	})
}

func (s *Service) GetRoleByCode(
	ctx context.Context,
	organizationID pgtype.UUID,
	code string,
) (db.Role, error) {
	if s.queries == nil {
		return db.Role{}, fmt.Errorf("rbac service queries are not initialized")
	}

	role, err := s.queries.GetRoleByCode(ctx, db.GetRoleByCodeParams{
		OrganizationID: organizationID,
		Code:           code,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Role{}, ErrRoleNotFound
		}

		return db.Role{}, fmt.Errorf("get role by code: %w", err)
	}

	return role, nil
}

func (s *Service) AssignRoleToUser(
	ctx context.Context,
	userID pgtype.UUID,
	roleID pgtype.UUID,
) error {
	if s.queries == nil {
		return fmt.Errorf("rbac service queries are not initialized")
	}

	if !userID.Valid {
		return fmt.Errorf("user ID is required")
	}

	if !roleID.Valid {
		return fmt.Errorf("role ID is required")
	}

	if err := s.queries.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	}); err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}

	return nil
}

func (s *Service) UserHasPermission(
	ctx context.Context,
	userID pgtype.UUID,
	permissionCode string,
) (bool, error) {
	if s.queries == nil {
		return false, fmt.Errorf("rbac service queries are not initialized")
	}

	if !userID.Valid {
		return false, fmt.Errorf("user ID is required")
	}

	if permissionCode == "" {
		return false, fmt.Errorf("permission code is required")
	}

	roles, err := s.queries.ListUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("list user roles: %w", err)
	}

	for _, role := range roles {
		permissions, err := s.queries.ListRolePermissions(ctx, role.ID)
		if err != nil {
			return false, fmt.Errorf("list role permissions: %w", err)
		}

		for _, permission := range permissions {
			if permission.Code == permissionCode {
				return true, nil
			}
		}
	}

	return false, nil
}
