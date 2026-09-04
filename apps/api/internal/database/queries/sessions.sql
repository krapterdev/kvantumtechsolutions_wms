-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    token_hash,
    expires_at
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    user_id,
    token_hash,
    expires_at,
    created_at,
    revoked_at;

-- name: GetSessionByTokenHash :one
SELECT
    id,
    user_id,
    token_hash,
    expires_at,
    created_at,
    revoked_at
FROM sessions
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET revoked_at = NOW()
WHERE id = $1
  AND revoked_at IS NULL;