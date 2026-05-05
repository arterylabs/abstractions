-- Rename deleted_at to disabled_at and deleted_by to disabled_reason in workspaces table
-- Also change deleted_by column type from UUID back to TEXT
ALTER TABLE workspaces
ALTER COLUMN deleted_by TYPE TEXT;
ALTER TABLE workspaces
  RENAME COLUMN deleted_at TO disabled_at;
ALTER TABLE workspaces
  RENAME COLUMN deleted_by TO disabled_reason;