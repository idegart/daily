CREATE TABLE slack_reports
(
    id               SERIAL primary key,
    slack_channel_id varchar(60),
    date             date         not null default CURRENT_DATE,
    ts               varchar(255) not null,

    created_at       timestamp             default current_timestamp,
    updated_at       timestamp             default current_timestamp,

    unique (slack_channel_id, date)
)