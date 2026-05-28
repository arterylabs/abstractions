CREATE INDEX idx_subscriptions_workspace_id ON subscriptions (workspace_id);

CREATE INDEX idx_invitations_workspace_id ON invitations (workspace_id);
CREATE INDEX idx_invitations_email ON invitations (email);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_impersonated_by ON sessions (impersonated_by) WHERE impersonated_by IS NOT NULL;
