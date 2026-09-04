package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/krapterdev/kvantumtechsolutions_wms/apps/api/internal/database/generated"
)

const sessionDuration = 24 * time.Hour

type SessionManager struct {
	queries *db.Queries
}

func NewSessionManager(queries *db.Queries) *SessionManager {
	return &SessionManager{
		queries: queries,
	}
}

func (m *SessionManager) CreateSession(
	ctx context.Context,
	userID pgtype.UUID,
) (string, error) {
	if m.queries == nil {
		return "", fmt.Errorf("session manager queries are not initialized")
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}

	rawToken := base64.RawURLEncoding.EncodeToString(tokenBytes)

	tokenHash := sha256.Sum256([]byte(rawToken))
	tokenHashString := base64.RawURLEncoding.EncodeToString(tokenHash[:])

	sessionID := uuid.New()
	sessionUUID := pgtype.UUID{
		Bytes: sessionID,
		Valid: true,
	}

	expiresAt := time.Now().UTC().Add(sessionDuration)

	_, err := m.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:        sessionUUID,
		UserID:    userID,
		TokenHash: tokenHashString,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return rawToken, nil
}

func (m *SessionManager) GetSessionByToken(
	ctx context.Context,
	rawToken string,
) (db.Session, error) {
	if m.queries == nil {
		return db.Session{}, fmt.Errorf("session manager queries are not initialized")
	}

	if rawToken == "" {
		return db.Session{}, fmt.Errorf("session token is required")
	}

	tokenHash := sha256.Sum256([]byte(rawToken))
	tokenHashString := base64.RawURLEncoding.EncodeToString(tokenHash[:])

	session, err := m.queries.GetSessionByTokenHash(ctx, tokenHashString)
	if err != nil {
		return db.Session{}, fmt.Errorf("get session by token: %w", err)
	}

	return session, nil
}

func (m *SessionManager) RevokeSession(
	ctx context.Context,
	sessionID pgtype.UUID,
) error {
	if m.queries == nil {
		return fmt.Errorf("session manager queries are not initialized")
	}

	if !sessionID.Valid {
		return fmt.Errorf("session ID is required")
	}

	if err := m.queries.RevokeSession(ctx, sessionID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
