CREATE TABLE users
(
    id          SERIAL primary key,
    email       varchar(60) not null unique,
    airtable_id integer unique,
    slack_id    varchar(60) unique,
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
)