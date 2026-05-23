ALTER TABLE repository_files
ADD COLUMN IF NOT EXISTS last_commit_at TIMESTAMP;
