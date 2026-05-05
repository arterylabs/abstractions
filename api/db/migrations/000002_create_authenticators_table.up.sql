CREATE TABLE IF NOT EXISTS authenticators (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  type TEXT NOT NULL,
  token TEXT NOT NULL,
  refresh_token TEXT,
  scopes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by UUID
);