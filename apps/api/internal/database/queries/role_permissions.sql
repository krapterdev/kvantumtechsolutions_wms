-- name: CreateRolePermission :exec
INSERT INTO role_permissions (
    role_id,
    permission_id
)
VALUES (
    $1,
    $2
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- name: RemoveRolePermission :exec
DELETE FROM role_permissions
WHERE role_id = $1
  AND permission_id = $2;

-- name: ListRolePermissions :many
SELECT
    p.id,
    p.code,
    p.name,
    p.description,
    p.created_at
FROM role_permissions rp
JOIN permissions p ON p.id = rp.permission_id
WHERE rp.role_id = $1
ORDER BY p.code;