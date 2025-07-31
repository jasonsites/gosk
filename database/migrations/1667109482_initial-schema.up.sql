CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_catalog.pg_type WHERE typname = 'record_status' AND typcategory = 'E') THEN
        CREATE TYPE record_status AS ENUM ('active', 'archived', 'deleted');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS example_entity (
  id                uuid                                PRIMARY KEY DEFAULT gen_random_uuid(),
  title             text                                NOT NULL,
  description       text,

  status            record_status     NOT NULL    DEFAULT 'active',
  created_on        timestamptz       NOT NULL    DEFAULT (now() at time zone 'utc'),
  created_context   jsonb             NOT NULL    DEFAULT '{}'::jsonb,
  modified_on       timestamptz       NOT NULL    DEFAULT (now() at time zone 'utc'),
  modified_context  jsonb             NOT NULL    DEFAULT '{}'::jsonb
);

CREATE INDEX example_entity_status_idx ON example_entity (status);
