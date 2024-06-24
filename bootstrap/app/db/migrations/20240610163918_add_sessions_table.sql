-- +goose Up
create table if not exists sessions(
	id integer primary key,
	token string not null,
	user_id integer not null references users,
	ip_address text,
	user_agent text,
	expires_at datetime not null, 
	created_at datetime not null, 
    updated_at datetime not null, 
	deleted_at datetime 
);

-- +goose Down
drop table if exists sessions;
