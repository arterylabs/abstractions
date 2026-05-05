-- Add impersonation tracking fields to sessions table
ALTER TABLE sessions
ADD COLUMN impersonated_at TIMESTAMP;
ALTER TABLE sessions
ADD COLUMN impersonated_by UUID;