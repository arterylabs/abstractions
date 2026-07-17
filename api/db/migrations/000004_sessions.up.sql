CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  token TEXT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by UUID,
  user_agent TEXT,
  ip TEXT,
  impersonated_at TIMESTAMP,
  impersonated_by UUID
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_impersonated_by ON sessions (impersonated_by) WHERE impersonated_by IS NOT NULL;
