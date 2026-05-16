ALTER TABLE newsletter_topics ADD COLUMN position INTEGER NOT NULL DEFAULT 0;

UPDATE newsletter_topics SET position = subq.row_num
FROM (
  SELECT id, ROW_NUMBER() OVER (PARTITION BY repo_id ORDER BY created_at ASC, id ASC) - 1 AS row_num
  FROM newsletter_topics
  WHERE deleted_at IS NULL AND sent_at IS NULL
) subq
WHERE newsletter_topics.id = subq.id;
