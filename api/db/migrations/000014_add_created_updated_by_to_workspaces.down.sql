-- Remove created_by and updated_by columns from workspaces table
ALTER TABLE workspaces DROP COLUMN created_by,
  DROP COLUMN updated_by;