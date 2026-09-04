-- name: GetUserByEmail :one
SELECT
    id,
    organization_id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active,
    created_at,
    updated_at
FROM users
WHERE organization_id = $1
  AND email = $2
LIMIT 1;


-- name: CreateUser :one
INSERT INTO users (
    id,
    organization_id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING
    id,
    organization_id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active,
    created_at,
    updated_at;


-- name: GetUserByID :one
SELECT
    id,
    organization_id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active,
    created_at,
    updated_at
FROM users
WHERE id = $1
LIMIT 1;