CREATE TABLE users
(
    id          SERIAL primary key,
    email       varchar(255) not null unique,
    name        varchar(255) not null,
    airtable_id integer unique,
    slack_id    varchar(60) unique,
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
)