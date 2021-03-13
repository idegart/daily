CREATE TABLE projects
(
    id              SERIAL primary key,
    name            varchar(255) not null,
    airtable_id     varchar(60)        default null,
    slack_id        varchar(60) unique default null,
    is_infographics bool               default false,
    created_at      timestamp          default current_timestamp,
    updated_at      timestamp          default current_timestamp
);

CREATE TABLE project_users
(
    id         SERIAL primary key,
    project_id int REFERENCES projects (id),
    user_id    int REFERENCES users (id),
    unique (project_id, user_id)
);