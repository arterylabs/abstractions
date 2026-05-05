-- Remove impersonation tracking fields from sessions table
ALTER TABLE sessions DROP COLUMN impersonated_by;
ALTER TABLE sessions DROP COLUMN impersonated_at;