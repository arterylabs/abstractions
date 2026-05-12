DROP INDEX IF EXISTS idx_integrations_workspace_type_resource;
ALTER TABLE integrations DROP COLUMN IF EXISTS resource_id;
ALTER TABLE integrations DROP COLUMN IF EXISTS resource_type;
