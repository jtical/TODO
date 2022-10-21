-- Filename: migrations/000001_add_list_indexes.up.sql
CREATE INDEX IF NOT EXISTS lists_name_idx ON lists USING GIN(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS lists_status_idx ON lists USING GIN(to_tsvector('simple', status));