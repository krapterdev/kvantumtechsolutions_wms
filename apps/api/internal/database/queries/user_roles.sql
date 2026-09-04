-- name: AssignRoleToUser :exec
INSERT INTO user_roles (
    user_id,
    role_id
)
VALUES (
    $1,
    $2
)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1
  AND role_id = $2;

-- name: ListUserRoles :many
SELECT
    r.id,
    r.organization_id,
    r.name,
    r.code,
    r.created_at
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
ORDER BY r.name;