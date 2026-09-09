ALTER TABLE dps.games
    ADD COLUMN slug                 TEXT UNIQUE,
    ADD COLUMN checksum             TEXT,
    ADD COLUMN igdb_updated_at      BIGINT,
    ADD COLUMN synced_at            TIMESTAMPTZ,
    RENAME COLUMN icon_link TO icon_url;

ALTER TABLE dps.platforms
    ADD COLUMN slug                 TEXT UNIQUE,
    ADD COLUMN icon_url             TEXT,
    ADD COLUMN checksum             TEXT,
    ADD COLUMN igdb_updated_at      BIGINT,
    ADD COLUMN synced_at            TIMESTAMPTZ;

ALTER TABLE dps.games     ALTER COLUMN id TYPE BIGINT;
ALTER TABLE dps.platforms ALTER COLUMN id TYPE BIGINT;
ALTER TABLE dps.teams     ALTER COLUMN game_id TYPE BIGINT;
ALTER TABLE dps.teams     ALTER COLUMN platform_id TYPE BIGINT;
ALTER TABLE dps.user_platforms ALTER COLUMN platform_id TYPE BIGINT;

-- Индекс для локального поиска по названию (ILIKE LIKE '%...%').
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_games_title_trgm     ON dps.games     USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_platforms_title_trgm ON dps.platforms USING gin (title gin_trgm_ops);
