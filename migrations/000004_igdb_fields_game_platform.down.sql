DROP INDEX IF EXISTS dps.idx_platforms_title_trgm;
DROP INDEX IF EXISTS dps.idx_games_title_trgm;

ALTER TABLE dps.platforms
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS icon_url,
    DROP COLUMN IF EXISTS checksum,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS synced_at;

ALTER TABLE dps.games
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS checksum,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS synced_at;

ALTER TABLE dps.games
    RENAME icon_url TO icon_link;

ALTER TABLE dps.games     ALTER COLUMN id TYPE INT;
ALTER TABLE dps.platforms ALTER COLUMN id TYPE INT;
ALTER TABLE dps.teams     ALTER COLUMN game_id TYPE INT;
ALTER TABLE dps.teams     ALTER COLUMN platform_id TYPE INT;
ALTER TABLE dps.user_platforms ALTER COLUMN platform_id TYPE INT;
