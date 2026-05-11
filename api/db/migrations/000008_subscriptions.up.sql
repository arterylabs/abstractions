CREATE TABLE IF NOT EXISTS subscriptions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id UUID NOT NULL,
  plan_id TEXT NOT NULL,
  status TEXT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  disabled_at TIMESTAMP,
  disabled_reason TEXT,
  "limit" INTEGER DEFAULT 1,
  deleted_at TIMESTAMPTZ,
  deleted_by UUID
);

CREATE INDEX idx_subscriptions_deleted_at ON subscriptions (deleted_at);
