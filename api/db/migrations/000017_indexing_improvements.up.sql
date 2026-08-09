-- Per-file processing failure isolation and indexer version invalidation.
ALTER TABLE repository_files
  ADD COLUMN IF NOT EXISTS process_error TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS process_attempts INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS indexer_version TEXT NOT NULL DEFAULT '';

-- Chunk provenance for citations and hybrid retrieval.
ALTER TABLE file_chunk_embeddings
  ADD COLUMN IF NOT EXISTS start_line INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS end_line INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS chunk_type TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS search_document TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS indexer_version TEXT NOT NULL DEFAULT '';

-- Full-text search support for hybrid (vector + keyword) retrieval.
ALTER TABLE file_chunk_embeddings
  ADD COLUMN IF NOT EXISTS search_tsv tsvector
  GENERATED ALWAYS AS (
    to_tsvector(
      'english',
      coalesce(search_document, '') || ' ' ||
      coalesce(summary, '') || ' ' ||
      coalesce(symbol_name, '') || ' ' ||
      coalesce(repo_file_path, '')
    )
  ) STORED;

CREATE INDEX IF NOT EXISTS file_chunk_embeddings_search_tsv_idx
  ON file_chunk_embeddings USING GIN (search_tsv);

CREATE INDEX IF NOT EXISTS file_chunk_embeddings_symbol_name_idx
  ON file_chunk_embeddings (repo_id, symbol_name)
  WHERE symbol_name <> '';
