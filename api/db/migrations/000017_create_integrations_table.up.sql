CREATE TABLE IF NOT EXISTS integrations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "type" TEXT NOT NULL,
  service_title TEXT,
  service_icon TEXT,
  workspace_id UUID NOT NULL REFERENCES workspaces(id),
  token TEXT NOT NULL,
  refresh_token TEXT,
  token_expires_at TIMESTAMP,
  scopes TEXT,
  service_tenant_id TEXT,
  service_user_id TEXT,
  "configuration" JSONB,
  created_by UUID NOT NULL REFERENCES users(id),
  updated_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  deleted_by UUID REFERENCES users(id)
);
-- Create indexes for better query performance
CREATE INDEX idx_integrations_workspace_id ON integrations(workspace_id);
CREATE INDEX idx_integrations_type ON integrations("type");
CREATE INDEX idx_integrations_deleted_at ON integrations(deleted_at);
CREATE INDEX idx_integrations_created_by ON integrations(created_by);