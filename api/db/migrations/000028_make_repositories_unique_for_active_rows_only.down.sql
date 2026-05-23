DROP INDEX IF EXISTS idx_repositories_workspace_service_active;

ALTER TABLE repositories
  ADD CONSTRAINT repositories_workspace_id_service_type_service_id_key
  UNIQUE (workspace_id, service_type, service_id);
