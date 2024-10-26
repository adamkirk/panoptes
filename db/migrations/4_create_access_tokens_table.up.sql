CREATE TABLE IF NOT EXISTS "user_access_tokens"(
   "id" TEXT PRIMARY KEY,
   "secret" TEXT NOT NULL,
   "expires_at" TIMESTAMP (6) WITH TIME ZONE,
   "user_id" UUID NOT NULL,
   CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS "user_access_tokens_id_expires_at_idx" ON "user_access_tokens" ("id", "expires_at");
COMMENT ON INDEX "user_access_tokens_id_expires_at_idx" IS 'id first, as it should have much higher cardinality. We''ll mostly query using both of these fields.';

COMMENT ON COLUMN "user_access_tokens"."id" IS 'The tokens ID, will be generated in a specific format by the app.';
COMMENT ON COLUMN "user_access_tokens"."secret" IS 'The secret part of the token, MUST be encrypted application-side before being stored.';
COMMENT ON COLUMN "user_access_tokens"."expires_at" IS 'When the token becomes invalid and stops working.';
COMMENT ON COLUMN "user_access_tokens"."user_id" IS 'A JSON blob containing all the relevant scopes for this user.';
