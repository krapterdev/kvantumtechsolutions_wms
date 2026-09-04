-- RBAC permissions
INSERT INTO permissions (id, code, name, description)
VALUES
    (gen_random_uuid(), 'users.read', 'View Users', 'View users'),
    (gen_random_uuid(), 'users.create', 'Create Users', 'Create new users'),
    (gen_random_uuid(), 'users.update', 'Update Users', 'Update existing users'),
    (gen_random_uuid(), 'users.delete', 'Delete Users', 'Delete users'),
    (gen_random_uuid(), 'roles.read', 'View Roles', 'View roles'),
    (gen_random_uuid(), 'roles.manage', 'Manage Roles', 'Create and manage roles'),
    (gen_random_uuid(), 'inventory.read', 'View Inventory', 'View inventory'),
    (gen_random_uuid(), 'inventory.manage', 'Manage Inventory', 'Create and update inventory')
ON CONFLICT (code) DO NOTHING;

-- Owner role
INSERT INTO roles (id, organization_id, name, code)
SELECT gen_random_uuid(), o.id, 'Owner', 'owner'
FROM organizations o
WHERE NOT EXISTS (
    SELECT 1 FROM roles r
    WHERE r.organization_id = o.id
      AND r.code = 'owner'
);

-- Administrator role
INSERT INTO roles (id, organization_id, name, code)
SELECT gen_random_uuid(), o.id, 'Administrator', 'admin'
FROM organizations o
WHERE NOT EXISTS (
    SELECT 1 FROM roles r
    WHERE r.organization_id = o.id
      AND r.code = 'admin'
);

-- Manager role
INSERT INTO roles (id, organization_id, name, code)
SELECT gen_random_uuid(), o.id, 'Manager', 'manager'
FROM organizations o
WHERE NOT EXISTS (
    SELECT 1 FROM roles r
    WHERE r.organization_id = o.id
      AND r.code = 'manager'
);

-- Owner gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'owner'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Administrator gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Manager permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'users.read',
    'roles.read',
    'inventory.read',
    'inventory.manage'
)
WHERE r.code = 'manager'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign Owner role to development admin user
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r
    ON r.organization_id = u.organization_id
   AND r.code = 'owner'
WHERE u.email = 'admin@kvantumtechsolutions.com'
ON CONFLICT (user_id, role_id) DO NOTHING;