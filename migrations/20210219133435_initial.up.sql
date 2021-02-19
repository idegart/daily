CREATE TABLE users
(
    id          SERIAL primary key,
    email       varchar(255) not null unique,
    name        varchar(255) not null,
    airtable_id integer unique default null,
    slack_id    varchar(60) unique default null,
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp
);

CREATE TABLE daily_user_reports
(
    id         SERIAL primary key,
    user_id    integer      not null references users (id),
    date       date         not null default CURRENT_DATE,
    done       varchar(255) not null,
    will_do    varchar(255) not null,
    blocker    varchar(255) not null,
    created_at timestamp             default current_timestamp,
    updated_at timestamp             default current_timestamp,

    unique (user_id, date)
);

CREATE TABLE daily_slack_reports
(
    id               SERIAL primary key,
    slack_channel_id varchar(60),
    date             date         not null default CURRENT_DATE,
    ts               varchar(255) not null,

    created_at       timestamp             default current_timestamp,
    updated_at       timestamp             default current_timestamp,

    unique (slack_channel_id, date)
)
