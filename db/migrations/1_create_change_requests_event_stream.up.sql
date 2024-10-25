CREATE TABLE IF NOT EXISTS "change_requests_stream"(
   "id" UUID PRIMARY KEY,
   "aggregate_id" TEXT NOT NULL,
   "occurred_at" TIMESTAMP (6) WITH TIME ZONE,
   "payload" TEXT NOT NULL,
   "type" TEXT,
   "source_id" UUID DEFAULT NULL,
   "source_integration" TEXT NOT NULL
);

COMMENT ON COLUMN "change_requests_stream"."aggregate_id" IS 'The id of the "aggregate" in the system.
This will be the ID of the change request from the source system. E.g. 1234 for github PR.';
COMMENT ON COLUMN "change_requests_stream"."occurred_at" IS 'Used for the projection, this controls where it sits in the timeline.';
COMMENT ON COLUMN "change_requests_stream"."payload" IS 'The actual payload for the event.
TEXT rather than a JSON type because we don''t wanna interact with it at a db level, thats what the projection and read side are for';
COMMENT ON COLUMN "change_requests_stream"."type" IS 'The type of event that occurred.';
COMMENT ON COLUMN "change_requests_stream"."source_id" IS 'The source id of the webhook that produced this event (if applicable).';
COMMENT ON COLUMN "change_requests_stream"."source_integration" IS 'The source integration name this will be something like ''github'', ''gitlab'', ''bitbucket'' etc
Mostly just meta so we can see where things are coming from and couples with source_id to create polymorphic key.
Hate polymorhpism in a database usually, but it makes sense here, as it shouldn''t be heavily used';


CREATE INDEX IF NOT EXISTS "aggregate_id_idx" ON "change_requests_stream" ("aggregate_id");
CREATE INDEX IF NOT EXISTS "occurred_at_idx" ON "change_requests_stream" ("occurred_at");
CREATE INDEX IF NOT EXISTS "type_idx" ON "change_requests_stream" ("type");

CREATE INDEX IF NOT EXISTS  "source_relation_idx" ON "change_requests_stream" ("source_id", "source_integration");
COMMENT ON INDEX "source_relation_idx" IS 'source_id first, as it should have much higher cardinality.';

CREATE INDEX IF NOT EXISTS  "type_idx" ON "change_requests_stream" ("aggregate_id", "type");
COMMENT ON INDEX "type_idx" IS 'aggregate_id first, as it should have much higher cardinality';