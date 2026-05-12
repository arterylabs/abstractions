ALTER TABLE integrations ADD COLUMN IF NOT EXISTS resource_type TEXT;
ALTER TABLE integrations ADD COLUMN IF NOT EXISTS resource_id UUID;

CREATE INDEX IF NOT EXISTS idx_integrations_workspace_type_resource
    ON integrations (workspace_id, type, resource_type, resource_id)
    WHERE deleted_at IS NULL;
