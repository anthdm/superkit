-- +goose Up
create table if not exists sessions(
	id integer primary key,
	token string not null,
	user_id integer not null references users,
	ip_address text,
	user_agent text,
	expires_at timestamp with time zone not null,
	last_login_at timestamp with time zone,
	created_at timestamp with time zone not null
);

-- +goose Down
drop table if exists sessions;
