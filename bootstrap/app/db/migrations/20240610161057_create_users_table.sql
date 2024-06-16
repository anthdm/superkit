-- +goose Up
create table if not exists users(
	id integer primary key,
	email text unique not null,
	password_hash text not null,
	first_name text not null,
	last_name text not null,
	email_verified_at timestamp with time zone,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone not null
);

-- +goose Down
drop table if exists users;
