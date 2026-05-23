ALTER TABLE repositories
  DROP CONSTRAINT IF EXISTS repositories_workspace_id_service_type_service_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_repositories_workspace_service_active
  ON repositories (workspace_id, service_type, service_id)
  WHERE deleted_at IS NULL;
