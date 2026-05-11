CREATE TABLE IF NOT EXISTS file_chunk_embeddings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  repo_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
  repo_file_path TEXT NOT NULL,
  repo_file_id UUID REFERENCES repository_files(id) ON DELETE
  SET NULL,
    chunk_index INT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    tags TEXT [] NOT NULL DEFAULT '{}',
    embedding VECTOR(1536),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX file_chunk_embeddings_repo_file_chunk_idx ON file_chunk_embeddings (repo_id, repo_file_path, chunk_index);
CREATE INDEX file_chunk_embeddings_repo_id_idx ON file_chunk_embeddings (repo_id);
CREATE INDEX file_chunk_embeddings_embedding_idx ON file_chunk_embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX idx_file_chunk_embeddings_deleted_at ON file_chunk_embeddings (deleted_at);