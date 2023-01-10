CREATE TABLE IF NOT EXISTS game_jams (
    game_jam_id INT GENERATED ALWAYS AS IDENTITY,
    name text,
    url text,
    content text,
    start_date timestamptz,
    end_date timestamptz,
    voting_end_date timestamptz,
    hide_results bool,
    hide_submissions bool,
    created_at timestamptz not null default current_timestamp);

CREATE TABLE IF NOT EXISTS games (
    game_id INT GENERATED ALWAYS AS IDENTITY,
    game_jam_id INT,
    name text,
    url text,
    content text,
    build text,
    is_banned bool,
    created_at timestamptz not null default current_timestamp);

---- create above / drop below ----

DROP TABLE IF EXISTS game_jams;
DROP TABLE IF EXISTS games;