ALTER TABLE dps.teams DROP CONSTRAINT teams_check;

ALTER TABLE dps.teams
    ADD CONSTRAINT teams_check
    CHECK (
        (is_rating_required = FALSE AND (desired_rating IS NULL OR desired_rating = ''))
        OR
        (is_rating_required = TRUE AND desired_rating IS NOT NULL AND desired_rating <> '')
    );

ALTER TABLE dps.teams
    ADD CONSTRAINT chk_teams_desired_rating_len
    CHECK (char_length(desired_rating) <= 200);

ALTER TABLE dps.teams
    ADD CONSTRAINT chk_teams_contact_link_len
    CHECK (char_length(contact_link) <= 2000);
