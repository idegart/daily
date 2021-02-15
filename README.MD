## Сборка
`go build -v -o bin/daily ./cmd/daily` - соберет проект в папку `./bin` 
## SQL
```
create table if not exists schema_migrations
(
	version bigint not null
		constraint schema_migrations_pkey
			primary key,
	dirty boolean not null
);

alter table schema_migrations owner to "slack-bot";

create table if not exists users
(
	id serial not null
		constraint users_pkey
			primary key,
	email varchar(255) not null
		constraint users_email_key
			unique,
	name varchar(255) not null,
	airtable_id integer
		constraint users_airtable_id_key
			unique,
	slack_id varchar(60)
		constraint users_slack_id_key
			unique,
	created_at timestamp default CURRENT_TIMESTAMP,
	updated_at timestamp default CURRENT_TIMESTAMP
);

alter table users owner to "slack-bot";

create table if not exists daily_reports
(
	id serial not null
		constraint daily_reports_pkey
			primary key,
	user_id integer not null
		constraint daily_reports_user_id_fkey
			references users,
	date date default CURRENT_DATE not null,
	done varchar(255) not null,
	will_do varchar(255) not null,
	blocker varchar(255) not null,
	created_at timestamp default CURRENT_TIMESTAMP,
	updated_at timestamp default CURRENT_TIMESTAMP,
	constraint daily_reports_user_id_date_key
		unique (user_id, date)
);

alter table daily_reports owner to "slack-bot";

create table if not exists slack_reports
(
	id serial not null
		constraint slack_reports_pkey
			primary key,
	slack_channel_id varchar(60),
	date date default CURRENT_DATE not null,
	ts varchar(255) not null,
	created_at timestamp default CURRENT_TIMESTAMP,
	updated_at timestamp default CURRENT_TIMESTAMP,
	constraint slack_reports_slack_channel_id_date_key
		unique (slack_channel_id, date)
);

alter table slack_reports owner to "slack-bot";


```