ALTER TABLE dps.teams DROP CONSTRAINT chk_teams_desired_rating_len;
ALTER TABLE dps.teams DROP CONSTRAINT chk_teams_contact_link_len;

ALTER TABLE dps.teams DROP CONSTRAINT teams_check;

ALTER TABLE dps.teams
    ADD CONSTRAINT teams_check
    CHECK (
        (is_rating_required = FALSE AND desired_rating IS NULL)
        OR
        (is_rating_required = TRUE AND desired_rating IS NOT NULL)
    );
