CREATE TABLE IF NOT EXISTS chunk_embeddings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id         UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    file_path       TEXT NOT NULL,
    chunk_index     INT NOT NULL,
    summary         TEXT NOT NULL DEFAULT '',
    domain_tags     TEXT[] NOT NULL DEFAULT '{}',
    embedding       VECTOR(1536),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS chunk_embeddings_repo_file_chunk_idx
    ON chunk_embeddings (repo_id, file_path, chunk_index);

CREATE INDEX IF NOT EXISTS chunk_embeddings_repo_id_idx
    ON chunk_embeddings (repo_id);

CREATE INDEX IF NOT EXISTS chunk_embeddings_embedding_idx
    ON chunk_embeddings USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);
