CREATE TABLE IF NOT EXISTS users
(
    user_id    INT GENERATED ALWAYS AS IDENTITY,
    email      text,
    username   text,
    password   bytea,
    avatar     text,
    about      text,
    created_at timestamptz not null default current_timestamp,

    primary key (user_id),
    unique (user_id)
);


CREATE TABLE IF NOT EXISTS game_jams
(
    game_jam_id      INT GENERATED ALWAYS AS IDENTITY,
    user_id          INT         NOT NULL REFERENCES users (user_id),
    title            text,
    url              text,
    content          text,
    cover_image      text,
    start_date       timestamptz,
    end_date         timestamptz,
    voting_end_date  timestamptz,
    hide_results     bool,
    hide_submissions bool,
    created_at       timestamptz not null default current_timestamp,

    primary key (game_jam_id),
    unique (game_jam_id)
);

CREATE TABLE IF NOT EXISTS games
(
    game_id           INT GENERATED ALWAYS AS IDENTITY,
    game_jam_id       INT         NOT NULL REFERENCES game_jams (game_jam_id),
    user_id           INT         NOT NULL REFERENCES users (user_id),
    title             text,
    url               text,
    content           text,
    cover_image       text,
    screenshot_images text[],
    build             text,
    is_banned         bool,
    created_at        timestamptz not null default current_timestamp,

    primary key (game_id),
    unique (game_id)
);

CREATE TABLE IF NOT EXISTS participants
(
    participant_id      INT GENERATED ALWAYS AS IDENTITY,
    user_id             INT         NOT NULL REFERENCES users (user_id),
    game_jam_id         INT         NOT NULL REFERENCES game_jams (game_jam_id),
    team_id             INT,
    is_looking_for_team bool,
    tags                text[],
    is_admin            bool,
    created_at          timestamptz not null default current_timestamp
);

CREATE TABLE IF NOT EXISTS sessions
(
    session_id text PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users (user_id),
    expire_at  timestamptz
);

CREATE TABLE IF NOT EXISTS criteria
(
    criteria_id INT GENERATED ALWAYS AS IDENTITY,
    jam_id INT NOT NULL REFERENCES game_jams (game_jam_id),
    title       text,
    description text,
    created_at  timestamptz not null default current_timestamp,

    primary key (criteria_id),
    unique (criteria_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS game_jams;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS criteria;
