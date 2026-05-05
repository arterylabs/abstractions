-- Add created_by and updated_by columns to workspaces table
ALTER TABLE workspaces
ADD COLUMN created_by UUID NOT NULL,
  ADD COLUMN updated_by UUID NOT NULL;