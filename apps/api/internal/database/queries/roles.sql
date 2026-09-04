-- name: CreateRole :one
INSERT INTO roles (
    id,
    organization_id,
    name,
    code
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    organization_id,
    name,
    code,
    created_at;

-- name: GetRoleByCode :one
SELECT
    id,
    organization_id,
    name,
    code,
    created_at
FROM roles
WHERE organization_id = $1
  AND code = $2
LIMIT 1;

-- name: ListRolesByOrganization :many
SELECT
    id,
    organization_id,
    name,
    code,
    created_at
FROM roles
WHERE organization_id = $1
ORDER BY name;