CREATE TABLE IF NOT EXISTS "github_webhooks"(
   "id" UUID PRIMARY KEY,
   "occurred_at" TIMESTAMP (6) WITH TIME ZONE,
   "payload" TEXT NOT NULL
);

COMMENT ON COLUMN "change_requests_stream"."occurred_at" IS 'Used for the projection, this controls where it sits in the timeline.';
COMMENT ON COLUMN "change_requests_stream"."payload" IS 'The actual payload for the event.
TEXT rather than a JSON type because we don''t wanna interact with it at a db level, just store and read it';
