DROP INDEX IF EXISTS idx_newsletter_topics_repo_id;
ALTER TABLE newsletter_topics DROP COLUMN IF EXISTS sent_at;
ALTER TABLE newsletter_topics DROP COLUMN IF EXISTS repo_id;
