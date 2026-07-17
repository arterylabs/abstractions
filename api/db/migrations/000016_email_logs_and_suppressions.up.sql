-- email_logs records every email send attempt (sent, failed, or suppressed) along with
-- the SES message id and any error, so we have a durable feedback loop for deliverability.
CREATE TABLE IF NOT EXISTS email_logs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT NOT NULL,
  message_id TEXT,
  subject TEXT,
  status TEXT NOT NULL,
  error_type TEXT,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_logs_email ON email_logs (email);
CREATE INDEX idx_email_logs_status ON email_logs (status);
CREATE INDEX idx_email_logs_created_at ON email_logs (created_at);

-- email_suppressions holds addresses we must not email: hard bounces, complaints,
-- repeated soft failures crossing the threshold, blocked domains, or manual entries.
-- Soft-deletable so an address can be manually unblocked while keeping history.
CREATE TABLE IF NOT EXISTS email_suppressions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT NOT NULL,
  reason TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX email_suppressions_email_idx ON email_suppressions (email) WHERE deleted_at IS NULL;
CREATE INDEX idx_email_suppressions_email ON email_suppressions (email);
