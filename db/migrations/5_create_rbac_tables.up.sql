CREATE TABLE IF NOT EXISTS "roles"(
   "id" UUID PRIMARY KEY,
   "name" TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "roles_unique_name" ON "roles" ("name");

CREATE TABLE IF NOT EXISTS "permissions"(
   "id" UUID PRIMARY KEY,
   "name" TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "permissions_unique_name" ON "permissions" ("name");

CREATE TABLE IF NOT EXISTS "roles_permissions"(
   "id" UUID PRIMARY KEY,
   "role_id" UUID NOT NULL,
   "permission_id" UUID NOT NULL,
   CONSTRAINT fk_role_id FOREIGN KEY(role_id) REFERENCES roles(id),
   CONSTRAINT fk_permission_id FOREIGN KEY(permission_id) REFERENCES permissions(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS "roles_permissions_unique_combination" ON "roles_permissions" ("role_id", "permission_id");

CREATE TABLE IF NOT EXISTS "user_roles"(
   "id" UUID PRIMARY KEY,
   "role_id" UUID NOT NULL,
   "user_id" UUID NOT NULL,
   CONSTRAINT fk_role_id FOREIGN KEY(role_id) REFERENCES roles(id),
   CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS "user_roles_unique_combination" ON "user_roles" ("role_id", "user_id");



