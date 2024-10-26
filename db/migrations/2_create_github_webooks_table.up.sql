CREATE TABLE IF NOT EXISTS "github_webhooks"(
   "id" UUID PRIMARY KEY,
   "occurred_at" TIMESTAMP (6) WITH TIME ZONE,
   "payload" JSON NOT NULL
);

COMMENT ON COLUMN "change_requests_stream"."occurred_at" IS 'Used for the projection, this controls where it sits in the timeline.';
COMMENT ON COLUMN "change_requests_stream"."payload" IS 'The actual payload for the event.
JSON should be fine, we don''t need to parse it in the DB.';
