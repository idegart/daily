CREATE TABLE daily_reports
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
)