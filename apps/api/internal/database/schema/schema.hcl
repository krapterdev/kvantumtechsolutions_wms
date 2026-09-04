schema "public" {
}

table "organizations" {
  schema = schema.public

  column "id" {
    type = uuid
  }

  column "name" {
    type = varchar(150)
  }

  column "code" {
    type = varchar(50)
  }

  column "is_active" {
    type    = boolean
    default = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "organizations_code_key" {
    unique  = true
    columns = [column.code]
  }
}

table "branches" {
  schema = schema.public

  column "id" {
    type = uuid
  }

  column "organization_id" {
    type = uuid
  }

  column "name" {
    type = varchar(150)
  }

  column "code" {
    type = varchar(50)
  }

  column "is_active" {
    type    = boolean
    default = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "branches_organization_code_key" {
    unique  = true
    columns = [column.organization_id, column.code]
  }

  foreign_key "branches_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organizations.column.id]
  }
}

table "roles" {
  schema = schema.public

  column "id" {
    type = uuid
  }

  column "organization_id" {
    type = uuid
  }

  column "name" {
    type = varchar(100)
  }

  column "code" {
    type = varchar(50)
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "roles_organization_code_key" {
    unique  = true
    columns = [column.organization_id, column.code]
  }

  foreign_key "roles_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organizations.column.id]
  }
}

table "users" {
  schema = schema.public

  column "id" {
    type = uuid
  }

  column "organization_id" {
    type = uuid
  }

  column "email" {
    type = varchar(255)
  }

  column "password_hash" {
    type = varchar(255)
  }

  column "first_name" {
    type = varchar(100)
  }

  column "last_name" {
    type = varchar(100)
    null = true
  }

  column "is_active" {
    type    = boolean
    default = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "users_organization_email_key" {
    unique  = true
    columns = [column.organization_id, column.email]
  }

  foreign_key "users_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organizations.column.id]
  }
}

table "user_branches" {
  schema = schema.public

  column "user_id" {
    type = uuid
  }

  column "branch_id" {
    type = uuid
  }

  primary_key {
    columns = [column.user_id, column.branch_id]
  }

  foreign_key "user_branches_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "user_branches_branch_id_fkey" {
    columns     = [column.branch_id]
    ref_columns = [table.branches.column.id]
    on_delete   = CASCADE
  }
}

table "user_roles" {
  schema = schema.public

  column "user_id" {
    type = uuid
  }

  column "role_id" {
    type = uuid
  }

  primary_key {
    columns = [column.user_id, column.role_id]
  }

  foreign_key "user_roles_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  foreign_key "user_roles_role_id_fkey" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_delete   = CASCADE
  }
}

table "permissions" {
  schema = schema.public

  column "id" {
    type = uuid
  }

  column "code" {
    type = varchar(100)
  }

  column "name" {
    type = varchar(150)
  }

  column "description" {
    type = varchar(255)
    null = true
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "permissions_code_key" {
    unique  = true
    columns = [column.code]
  }
}

table "role_permissions" {
  schema = schema.public

  column "role_id" {
    type = uuid
  }

  column "permission_id" {
    type = uuid
  }

  primary_key {
    columns = [column.role_id, column.permission_id]
  }

  foreign_key "role_permissions_role_id_fkey" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_delete   = CASCADE
  }

  foreign_key "role_permissions_permission_id_fkey" {
    columns     = [column.permission_id]
    ref_columns = [table.permissions.column.id]
    on_delete   = CASCADE
  }
}

table "sessions" {
  schema = schema.public

  column "id" {
    type = uuid
  }

  column "user_id" {
    type = uuid
  }

  column "token_hash" {
    type = varchar(255)
  }

  column "expires_at" {
    type = timestamptz
  }

  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  column "revoked_at" {
    type = timestamptz
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  index "sessions_token_hash_key" {
    unique  = true
    columns = [column.token_hash]
  }

  index "sessions_user_id_idx" {
    columns = [column.user_id]
  }

  index "sessions_expires_at_idx" {
    columns = [column.expires_at]
  }

  foreign_key "sessions_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
}