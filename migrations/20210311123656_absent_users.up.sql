CREATE TABLE absent_users
(
    id      SERIAL primary key,
    user_id int REFERENCES users (id),
    date    date,
    unique (user_id, date)
);