CREATE TABLE IF NOT EXISTS repositories (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id UUID NOT NULL REFERENCES workspaces(id),
  user_id UUID NOT NULL REFERENCES users(id),
  github_repo_id BIGINT NOT NULL,
  full_name TEXT NOT NULL,
  default_branch TEXT NOT NULL DEFAULT 'main',
  webhook_id BIGINT,
  webhook_secret TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  indexed_at TIMESTAMP,
  send_day INTEGER NOT NULL DEFAULT 1,
  send_hour INTEGER NOT NULL DEFAULT 9,
  last_newsletter_sent_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  UNIQUE (workspace_id, github_repo_id)
);

CREATE INDEX idx_repositories_workspace_id ON repositories(workspace_id);
CREATE INDEX idx_repositories_user_id ON repositories(user_id);
CREATE INDEX idx_repositories_status ON repositories(status);
CREATE INDEX idx_repositories_deleted_at ON repositories(deleted_at);
CREATE INDEX idx_repositories_send_day_hour ON repositories(send_day, send_hour);
