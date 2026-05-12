ALTER TABLE newsletter_topics ADD COLUMN repo_id UUID REFERENCES repositories(id) ON DELETE CASCADE;
ALTER TABLE newsletter_topics ADD COLUMN sent_at TIMESTAMPTZ;
CREATE INDEX idx_newsletter_topics_repo_id ON newsletter_topics(repo_id);
