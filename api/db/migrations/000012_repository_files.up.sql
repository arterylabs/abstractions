CREATE TABLE IF NOT EXISTS repository_files (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  repo_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
  path TEXT NOT NULL,
  git_sha TEXT NOT NULL DEFAULT '',
  tags TEXT [] NOT NULL DEFAULT '{}',
  processed_at TIMESTAMP,
  size INTEGER NOT NULL DEFAULT 0,
  language TEXT NOT NULL DEFAULT '',
  last_commit_at TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  deleted_by UUID,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_repository_files_repo_path ON repository_files (repo_id, path);
CREATE INDEX idx_repository_files_git_sha ON repository_files (git_sha);
CREATE INDEX idx_repository_files_deleted_at ON repository_files (deleted_at);