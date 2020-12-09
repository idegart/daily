CREATE TABLE users
(
    id          SERIAL primary key,
    airtable_id integer unique,
    slack_id    varchar(60) unique,
    name        varchar(60) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);