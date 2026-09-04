-- Create "organizations" table
CREATE TABLE "public"."organizations" (
  "id" uuid NOT NULL,
  "name" character varying(150) NOT NULL,
  "code" character varying(50) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "organizations_code_key" to table: "organizations"
CREATE UNIQUE INDEX "organizations_code_key" ON "public"."organizations" ("code");
-- Create "branches" table
CREATE TABLE "public"."branches" (
  "id" uuid NOT NULL,
  "organization_id" uuid NOT NULL,
  "name" character varying(150) NOT NULL,
  "code" character varying(50) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "branches_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "branches_organization_code_key" to table: "branches"
CREATE UNIQUE INDEX "branches_organization_code_key" ON "public"."branches" ("organization_id", "code");
-- Create "permissions" table
CREATE TABLE "public"."permissions" (
  "id" uuid NOT NULL,
  "code" character varying(100) NOT NULL,
  "name" character varying(150) NOT NULL,
  "description" character varying(255) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "permissions_code_key" to table: "permissions"
CREATE UNIQUE INDEX "permissions_code_key" ON "public"."permissions" ("code");
-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL,
  "organization_id" uuid NOT NULL,
  "name" character varying(100) NOT NULL,
  "code" character varying(50) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "roles_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "roles_organization_code_key" to table: "roles"
CREATE UNIQUE INDEX "roles_organization_code_key" ON "public"."roles" ("organization_id", "code");
-- Create "role_permissions" table
CREATE TABLE "public"."role_permissions" (
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("role_id", "permission_id"),
  CONSTRAINT "role_permissions_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "organization_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" character varying(255) NOT NULL,
  "first_name" character varying(100) NOT NULL,
  "last_name" character varying(100) NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "users_organization_email_key" to table: "users"
CREATE UNIQUE INDEX "users_organization_email_key" ON "public"."users" ("organization_id", "email");
-- Create "sessions" table
CREATE TABLE "public"."sessions" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "token_hash" character varying(255) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "revoked_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "sessions_expires_at_idx" to table: "sessions"
CREATE INDEX "sessions_expires_at_idx" ON "public"."sessions" ("expires_at");
-- Create index "sessions_token_hash_key" to table: "sessions"
CREATE UNIQUE INDEX "sessions_token_hash_key" ON "public"."sessions" ("token_hash");
-- Create index "sessions_user_id_idx" to table: "sessions"
CREATE INDEX "sessions_user_id_idx" ON "public"."sessions" ("user_id");
-- Create "user_branches" table
CREATE TABLE "public"."user_branches" (
  "user_id" uuid NOT NULL,
  "branch_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "branch_id"),
  CONSTRAINT "user_branches_branch_id_fkey" FOREIGN KEY ("branch_id") REFERENCES "public"."branches" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_branches_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "user_roles" table
CREATE TABLE "public"."user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "role_id"),
  CONSTRAINT "user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_roles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
