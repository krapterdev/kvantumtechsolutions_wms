-- name: CreatePermission :one
INSERT INTO permissions (
    id,
    code,
    name,
    description
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    code,
    name,
    description,
    created_at;

-- name: GetPermissionByCode :one
SELECT
    id,
    code,
    name,
    description,
    created_at
FROM permissions
WHERE code = $1
LIMIT 1;

-- name: ListPermissions :many
SELECT
    id,
    code,
    name,
    description,
    created_at
FROM permissions
ORDER BY code;