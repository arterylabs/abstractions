CREATE TABLE IF NOT EXISTS invitations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id UUID NOT NULL,
  invited_by UUID NOT NULL,
  email TEXT NOT NULL,
  role TEXT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  accepted_at TIMESTAMP,
  ignored_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  deleted_by UUID
);

CREATE INDEX idx_invitations_deleted_at ON invitations (deleted_at);
CREATE INDEX idx_invitations_workspace_id ON invitations (workspace_id);
CREATE INDEX idx_invitations_email ON invitations (email);
