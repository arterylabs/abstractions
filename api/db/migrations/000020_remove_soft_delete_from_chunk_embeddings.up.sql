DROP INDEX IF EXISTS idx_file_chunk_embeddings_deleted_at;
ALTER TABLE file_chunk_embeddings DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE file_chunk_embeddings DROP COLUMN IF EXISTS deleted_by;
