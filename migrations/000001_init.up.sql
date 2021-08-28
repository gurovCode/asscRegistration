CREATE TABLE if not exists teams
(
    id                    bigserial primary key,
    created_at            timestamp with time zone default current_timestamp,
    updated_at            timestamp with time zone default current_timestamp,
    name                  text not null unique,
    region                text not null,
    university            text not null,
    faceit_link           text not null unique,
    is_confirmed          bool not null default false
);

CREATE TABLE if not exists players
(
    id                    bigserial primary key,
    created_at            timestamp with time zone default current_timestamp,
    updated_at            timestamp with time zone default current_timestamp,
    nickname              text not null unique,
    team_id               bigint not null references teams(id),
    is_captain            bool not null default false,
    first_name            text not null,
    second_name           text,
    last_name             text not null,
    birthdate             date not null,
    phone_number          text unique,
    email                 text unique,
    discord_id            text unique,
    document              text not null unique
);
