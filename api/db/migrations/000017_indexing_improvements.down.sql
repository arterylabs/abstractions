DROP INDEX IF EXISTS file_chunk_embeddings_symbol_name_idx;
DROP INDEX IF EXISTS file_chunk_embeddings_search_tsv_idx;

ALTER TABLE file_chunk_embeddings
  DROP COLUMN IF EXISTS search_tsv,
  DROP COLUMN IF EXISTS indexer_version,
  DROP COLUMN IF EXISTS search_document,
  DROP COLUMN IF EXISTS language,
  DROP COLUMN IF EXISTS chunk_type,
  DROP COLUMN IF EXISTS end_line,
  DROP COLUMN IF EXISTS start_line;

ALTER TABLE repository_files
  DROP COLUMN IF EXISTS indexer_version,
  DROP COLUMN IF EXISTS process_attempts,
  DROP COLUMN IF EXISTS process_error;
