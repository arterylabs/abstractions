CREATE TABLE IF NOT EXISTS newsletters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    topics JSONB NOT NULL DEFAULT '[]',
    html_content TEXT NOT NULL DEFAULT '',
    sent_at TIMESTAMP,
    recipient_count INT NOT NULL DEFAULT 0,
    open_count INT NOT NULL DEFAULT 0,
    click_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_newsletters_repo_id ON newsletters (repo_id);
CREATE INDEX idx_newsletters_sent_at ON newsletters (sent_at);
