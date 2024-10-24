CREATE TABLE IF NOT EXISTS tasks_stream(
   id UUID PRIMARY KEY,

    -- the id of the "aggregate" in the system
    -- this will depend on whether the event is linked to a pull request for 
    -- example, where this will be an id for that, or it could be a commit sha etc
    -- don't think we ever want
   aggregate_id TEXT NOT NULL,

    -- used for the projection, this controls where it sits in the timeline
   occurred_at TIMESTAMP (6) WITH TIME ZONE,

    -- when the event was received from the source
    -- largely useful as meta, and shouldn't affect any projection
   received_at TIMESTAMP (6) WITH TIME ZONE,

   -- the actual payload for the event
   -- text because we don't wanna interact with it at a db level, thats what the
   -- projection and read side are for
   payload TEXT NOT NULL,

   -- the type of event
   type TEXT
);

CREATE INDEX IF NOT EXISTS  aggregate_id_idx ON tasks_stream (aggregate_id);
CREATE INDEX IF NOT EXISTS  occurred_at_idx ON tasks_stream (occurred_at);
CREATE INDEX IF NOT EXISTS  received_at_idx ON tasks_stream (received_at);
CREATE INDEX IF NOT EXISTS  type_idx ON tasks_stream (type);

-- aggregate id first as i expect it has higher cardinality than type
CREATE INDEX IF NOT EXISTS  type_idx ON tasks_stream (aggregate_id, type);