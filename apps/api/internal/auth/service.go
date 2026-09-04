package auth

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
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

type User struct {
	ID             pgtype.UUID
	OrganizationID pgtype.UUID
	Email          string
	PasswordHash   string
	FirstName      string
	LastName       pgtype.Text
	IsActive       bool
}

func (s *Service) GetUserByEmail(
	ctx context.Context,
	organizationID pgtype.UUID,
	email string,
) (User, error) {
	if s.queries == nil {
		return User{}, fmt.Errorf("auth service queries are not initialized")
	}

	user, err := s.queries.GetUserByEmail(ctx, db.GetUserByEmailParams{
		OrganizationID: organizationID,
		Email:          email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}

		return User{}, fmt.Errorf("get user by email: %w", err)
	}

	return User{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsActive:       user.IsActive,
	}, nil
}

func (s *Service) GetUserByID(
	ctx context.Context,
	userID pgtype.UUID,
) (User, error) {
	if s.queries == nil {
		return User{}, fmt.Errorf("auth service queries are not initialized")
	}

	if !userID.Valid {
		return User{}, fmt.Errorf("user ID is required")
	}

	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("user not found")
		}

		return User{}, fmt.Errorf("get user by ID: %w", err)
	}

	return User{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsActive:       user.IsActive,
	}, nil
}

type LoginRequest struct {
	OrganizationID pgtype.UUID
	Email          string
	Password       string
}

type LoginResult struct {
	User         User
	SessionToken string
}

func (s *Service) Login(
	ctx context.Context,
	sessionManager *SessionManager,
	req LoginRequest,
) (LoginResult, error) {
	if s.queries == nil {
		return LoginResult{}, fmt.Errorf("auth service queries are not initialized")
	}

	if sessionManager == nil {
		return LoginResult{}, fmt.Errorf("session manager is not initialized")
	}

	if !req.OrganizationID.Valid {
		return LoginResult{}, ErrInvalidCredentials
	}

	if req.Email == "" || req.Password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	user, err := s.GetUserByEmail(ctx, req.OrganizationID, req.Email)
	if err != nil {
		return LoginResult{}, err
	}

	if !user.IsActive {
		return LoginResult{}, ErrUserInactive
	}

	if !VerifyPassword(req.Password, user.PasswordHash) {
		return LoginResult{}, ErrInvalidCredentials
	}

	sessionToken, err := sessionManager.CreateSession(ctx, user.ID)
	if err != nil {
		return LoginResult{}, fmt.Errorf("create login session: %w", err)
	}

	return LoginResult{
		User:         user,
		SessionToken: sessionToken,
	}, nil
}

type CreateUserRequest struct {
	OrganizationID pgtype.UUID
	Email          string
	Password       string
	FirstName      string
	LastName       *string
}

func (s *Service) CreateUser(
	ctx context.Context,
	req CreateUserRequest,
) (User, error) {
	if s.queries == nil {
		return User{}, fmt.Errorf("auth service queries are not initialized")
	}

	if !req.OrganizationID.Valid {
		return User{}, fmt.Errorf("organization ID is required")
	}

	if req.Email == "" {
		return User{}, fmt.Errorf("email is required")
	}

	if req.Password == "" {
		return User{}, fmt.Errorf("password is required")
	}

	if req.FirstName == "" {
		return User{}, fmt.Errorf("first name is required")
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	userID := pgtype.UUID{
		Bytes: uuid.New(),
		Valid: true,
	}

	var lastName pgtype.Text
	if req.LastName != nil {
		lastName = pgtype.Text{
			String: *req.LastName,
			Valid:  true,
		}
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		ID:             userID,
		OrganizationID: req.OrganizationID,
		Email:          req.Email,
		PasswordHash:   passwordHash,
		FirstName:      req.FirstName,
		LastName:       lastName,
		IsActive:       true,
	})
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return User{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsActive:       user.IsActive,
	}, nil
}
