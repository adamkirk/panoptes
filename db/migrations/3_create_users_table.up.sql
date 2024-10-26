CREATE TABLE IF NOT EXISTS "users"(
   "id" UUID PRIMARY KEY,
   "email" TEXT NOT NULL,
   "first_name" TEXT NOT NULL,
   "last_name" TEXT NOT NULL,
   "password" TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_email_unique_idx" ON "users" ("email");

COMMENT ON COLUMN "users"."password" IS 'The users password (hashed), cannot be read at a database level.';
