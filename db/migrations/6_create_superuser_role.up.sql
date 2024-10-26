-- Don't wanna bother installing uuid extension in postgres, for now a static id
-- is probably fine...can't see any major issues.
INSERT INTO "roles" ("id", "name") VALUES ('7f137080-298f-461d-a3c7-764f911f59f5', 'superuser');