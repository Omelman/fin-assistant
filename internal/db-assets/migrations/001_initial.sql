-- +migrate Up

create table users (
  id bigserial primary key not null,
  firstname varchar(40) not null,
  lastname varchar(40) not null,
  email varchar(254) unique not null,
  password bytea not null,
  token varchar(64)
);

create table work (
  id bigserial primary key not null,
  date_start date not null default current_date,
  date_finish date,
  position varchar(100) not null,
  city varchar(20) not null,
  company varchar(30) not null,
  user_id bigserial references users (id) not null
);

create table balance (
  id bigserial primary key not null,
  currency varchar(4) not null,
  user_id bigserial references users (id) not null
);

create table transaction (
  id bigserial primary key not null,
  date date not null default current_date,
  description text not null,
  amount int not null,
  category varchar(20) not null,
  include boolean not null,
  balance_id bigserial not null references balance (id) on delete cascade on update cascade
);

create table goal (
  id bigserial primary key not null,
  date_start date not null default current_date,
  date_finish date,
  description text not null,
  amount int not null,
  balance_id bigserial references balance (id) not null
);

-- +migrate Down

drop table users;
drop table work;
drop table balance;
drop table transaction;
drop table goal;

