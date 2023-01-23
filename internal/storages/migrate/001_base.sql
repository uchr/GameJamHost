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

CREATE TABLE IF NOT EXISTS sessions
(
    session_id text NOT NULL,
    user_id    INT NOT NULL REFERENCES users (user_id),
    expire_at  timestamptz,

    primary key (session_id),
    unique (session_id)
);

CREATE TABLE IF NOT EXISTS criteria
(
    criteria_id int GENERATED ALWAYS AS IDENTITY,
    jam_id INT NOT NULL REFERENCES game_jams (game_jam_id),
    title       text,
    description text,
    created_at  timestamptz not null default current_timestamp,

    primary key (criteria_id),
    unique (criteria_id)
);

CREATE TABLE IF NOT EXISTS jam_questions
(
    question_id int GENERATED ALWAYS AS IDENTITY,
    jam_id INT NOT NULL REFERENCES game_jams (game_jam_id),
    title       text,
    description text,
    hidden_criteria text,
    created_at  timestamptz not null default current_timestamp,

    primary key (question_id),
    unique (question_id)
);

CREATE TABLE IF NOT EXISTS game_answers
(
    answer_id int GENERATED ALWAYS AS IDENTITY,
    game_id INT NOT NULL REFERENCES games (game_id),
    question_id int NOT NULL REFERENCES jam_questions (question_id),
    answer bool,

    primary key (answer_id),
    unique (answer_id)
);

CREATE TABLE IF NOT EXISTS votes
(
    vote_id int GENERATED ALWAYS AS IDENTITY,
    game_id INT NOT NULL REFERENCES games (game_id),
    user_id INT NOT NULL REFERENCES users (user_id),
    criteria_id int NOT NULL REFERENCES criteria (criteria_id),
    value int,

    primary key (vote_id),
    unique (vote_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS game_jams;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS criteria;
DROP TABLE IF EXISTS jam_questions;
DROP TABLE IF EXISTS game_answers;
DROP TABLE IF EXISTS votes;