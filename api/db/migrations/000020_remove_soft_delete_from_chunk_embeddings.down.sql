ALTER TABLE file_chunk_embeddings ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE file_chunk_embeddings ADD COLUMN IF NOT EXISTS deleted_by UUID;
CREATE INDEX IF NOT EXISTS idx_file_chunk_embeddings_deleted_at ON file_chunk_embeddings (deleted_at);
