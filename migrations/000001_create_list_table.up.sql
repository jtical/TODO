-- Filename: migrations/000001_create_list_table.up.sql

CREATE TABLE IF NOT EXISTS lists (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    task text NOT NULL,
    status text NOT NULL,
    version integer NOT NULL DEFAULT 1
);