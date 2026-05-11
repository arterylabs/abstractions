CREATE TABLE IF NOT EXISTS newsletters (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  repo_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
  topic_id UUID REFERENCES newsletter_topics(id) ON DELETE SET NULL,
  html_content TEXT NOT NULL DEFAULT '',
  sent_at TIMESTAMP,
  recipient_count INT NOT NULL DEFAULT 0,
  open_count INT NOT NULL DEFAULT 0,
  click_count INT NOT NULL DEFAULT 0,
  deleted_at TIMESTAMPTZ,
  deleted_by UUID,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_newsletters_repo_id ON newsletters (repo_id);
CREATE INDEX idx_newsletters_sent_at ON newsletters (sent_at);
CREATE INDEX idx_newsletters_deleted_at ON newsletters (deleted_at);
