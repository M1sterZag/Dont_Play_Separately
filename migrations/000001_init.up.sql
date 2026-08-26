CREATE SCHEMA dps;


CREATE TABLE dps.users (
    id              UUID         PRIMARY KEY,
    version         INTEGER      NOT NULL DEFAULT 1,
    email           VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    nickname        VARCHAR(40)  NOT NULL
        CHECK (char_length(nickname) BETWEEN 1 AND 40),
        
    bio             TEXT,
    avatar_url      TEXT,
    created_at      TIMESTAMPTZ  NOT NULL
);


CREATE TABLE dps.games (
    id        INTEGER PRIMARY KEY,
    title     TEXT    NOT NULL,
    icon_link TEXT
);


CREATE TABLE dps.platforms (
    id           INTEGER PRIMARY KEY,
    title        TEXT    NOT NULL,
    abbreviation TEXT    NOT NULL
);


CREATE TABLE dps.user_platforms (
    user_id     UUID    NOT NULL,
    platform_id INTEGER NOT NULL,

    PRIMARY KEY (user_id, platform_id)
);


CREATE TABLE dps.teams (
    id                 UUID         PRIMARY KEY,
    version            INTEGER      NOT NULL DEFAULT 1,
    owner_id           UUID         NOT NULL,
    game_id            INTEGER      NOT NULL,
    platform_id        INTEGER      NOT NULL,

    title              VARCHAR(100) NOT NULL
        CHECK (char_length(title) BETWEEN 1 AND 100),

    description        TEXT
        CHECK (char_length(description) BETWEEN 1 AND 1000),

    is_rating_required BOOLEAN      NOT NULL DEFAULT FALSE,
    desired_rating    TEXT,
    contact_link      TEXT,

    slots_total       INTEGER      NOT NULL
        CHECK (slots_total > 0),

    created_at         TIMESTAMPTZ  NOT NULL,
    is_active          BOOLEAN      NOT NULL DEFAULT TRUE,

    CHECK (
        (is_rating_required = FALSE AND desired_rating IS NULL)
        OR
        (is_rating_required = TRUE AND desired_rating IS NOT NULL)
    )
);


CREATE TABLE dps.team_members (
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,

    PRIMARY KEY (team_id, user_id)
);


ALTER TABLE dps.teams
    ADD CONSTRAINT fk_teams_owner
    FOREIGN KEY (owner_id)
    REFERENCES dps.users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;


ALTER TABLE dps.teams
    ADD CONSTRAINT fk_teams_game
    FOREIGN KEY (game_id)
    REFERENCES dps.games(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;


ALTER TABLE dps.teams
    ADD CONSTRAINT fk_teams_platform
    FOREIGN KEY (platform_id)
    REFERENCES dps.platforms(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;


ALTER TABLE dps.team_members
    ADD CONSTRAINT fk_team_members_team
    FOREIGN KEY (team_id)
    REFERENCES dps.teams(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;


ALTER TABLE dps.team_members
    ADD CONSTRAINT fk_team_members_user
    FOREIGN KEY (user_id)
    REFERENCES dps.users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;


ALTER TABLE dps.user_platforms
    ADD CONSTRAINT fk_user_platforms_user
    FOREIGN KEY (user_id)
    REFERENCES dps.users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;


ALTER TABLE dps.user_platforms
    ADD CONSTRAINT fk_user_platforms_platform
    FOREIGN KEY (platform_id)
    REFERENCES dps.platforms(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;