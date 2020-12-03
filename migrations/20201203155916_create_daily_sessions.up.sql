CREATE TABLE daily_session
(
    id   SERIAL primary key,
    date date not null default current_date,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);