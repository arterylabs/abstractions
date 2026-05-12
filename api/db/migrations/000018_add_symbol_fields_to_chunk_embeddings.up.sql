ALTER TABLE file_chunk_embeddings ADD COLUMN IF NOT EXISTS symbol_name TEXT NOT NULL DEFAULT '';
ALTER TABLE file_chunk_embeddings ADD COLUMN IF NOT EXISTS symbol_type TEXT NOT NULL DEFAULT '';
