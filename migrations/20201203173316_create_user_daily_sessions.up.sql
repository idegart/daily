CREATE TABLE user_daily_sessions
(
    id SERIAL primary key,
    user_id int not null references users,
    daily_session_id int not null references daily_session,
    done text default '',
    will_do text default '',
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
)