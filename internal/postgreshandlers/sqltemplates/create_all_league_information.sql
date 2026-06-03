CREATE TABLE IF NOT EXISTS all_league_information (
    league_key              TEXT        NOT NULL,
    league_id               INTEGER     NOT NULL,
    league_name             TEXT        NOT NULL,
    league_logo_url         TEXT,
    number_of_teams         INTEGER     NOT NULL,
    league_update_timestamp TIMESTAMPTZ,
    start_date              TEXT,
    end_week                TEXT,
    season                  INTEGER     NOT NULL,
    CONSTRAINT pk_all_league_information PRIMARY KEY (league_key)
);
