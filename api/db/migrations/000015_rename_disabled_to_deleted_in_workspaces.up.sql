-- Rename disabled_at to deleted_at and disabled_reason to deleted_by in workspaces table
-- Also change deleted_by column type from TEXT to UUID
ALTER TABLE workspaces
  RENAME COLUMN disabled_at TO deleted_at;
ALTER TABLE workspaces
  RENAME COLUMN disabled_reason TO deleted_by;
ALTER TABLE workspaces
ALTER COLUMN deleted_by TYPE UUID USING deleted_by::UUID;