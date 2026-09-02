CREATE TABLE dps.refresh_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE dps.refresh_sessions
    ADD CONSTRAINT fk_refresh_sessions_user
    FOREIGN KEY (user_id) REFERENCES dps.users (id)
    ON UPDATE CASCADE ON DELETE CASCADE;