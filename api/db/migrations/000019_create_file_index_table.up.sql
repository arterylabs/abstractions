CREATE TABLE IF NOT EXISTS file_index (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    git_sha TEXT NOT NULL DEFAULT '',
    domain_tags TEXT[] NOT NULL DEFAULT '{}',
    last_processed_at TIMESTAMP,
    pinecone_ids TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_file_index_repo_path ON file_index (repo_id, file_path);
CREATE INDEX idx_file_index_git_sha ON file_index (git_sha);
