CREATE TABLE workspace_configurations (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID        NOT NULL REFERENCES workspaces(id),
    key          TEXT        NOT NULL,
    value        TEXT        NOT NULL,
    encrypted    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by   UUID        NOT NULL,
    updated_by   UUID        NOT NULL,
    UNIQUE (workspace_id, key)
);
