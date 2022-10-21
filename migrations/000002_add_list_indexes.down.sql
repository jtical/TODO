-- Filename: migrations/000001_add_list_indexes.down.sql

DROP INDEX IF EXISTS lists_name_idx;
DROP INDEX IF EXISTS lists_status_idx;