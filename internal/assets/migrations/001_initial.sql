-- +migrate Up

create table users (
  id bigserial primary key not null,
  nickname varchar(40) not null,
  email varchar(254) unique not null,
  details jsonb,
  password bytea not null,
  token varchar(64),
  recovery_key varchar(64)
);

-- +migrate Down

drop table users;


