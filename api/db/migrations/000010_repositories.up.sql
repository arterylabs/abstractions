CREATE TABLE IF NOT EXISTS repositories (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id UUID NOT NULL REFERENCES workspaces(id),
  created_by UUID NOT NULL REFERENCES users(id),
  updated_by UUID REFERENCES users(id),
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  service_type TEXT NOT NULL DEFAULT 'GITHUB',
  service_id TEXT NOT NULL,
  default_branch TEXT NOT NULL DEFAULT 'main',
  status TEXT NOT NULL DEFAULT 'pending',
  indexed_at TIMESTAMP,
  index_error TEXT,
  send_day INTEGER NOT NULL DEFAULT 1,
  send_hour INTEGER NOT NULL DEFAULT 9,
  newsletter_sent_at TIMESTAMP,
  recipients TEXT[] NOT NULL DEFAULT '{}',
  send_to_workspace_members BOOLEAN NOT NULL DEFAULT TRUE,
  newsletter_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  deleted_by UUID
);

CREATE UNIQUE INDEX idx_repositories_workspace_service_active
  ON repositories (workspace_id, service_type, service_id)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_repositories_workspace_id ON repositories (workspace_id);
CREATE INDEX idx_repositories_created_by ON repositories (created_by);
CREATE INDEX idx_repositories_status ON repositories (status);
CREATE INDEX idx_repositories_deleted_at ON repositories (deleted_at);
CREATE INDEX idx_repositories_send_day_hour ON repositories (send_day, send_hour);
