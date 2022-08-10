CREATE TABLE if not exists users (
    id bigserial primary key,
    firstname varchar(15) not null UNIQUE,
    lastname varchar(30) not null,
    password varchar(100) not null,
    created_at timestamp default current_timestamp,
    last_seen_at timestamp default current_timestamp
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists todo(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(30) not null,
    deadline timestamp default current_timestamp,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    created_by bigserial,
    foreign key (created_by) references users(id)
    ON UPDATE CASCADE ON DELETE cascade,
    overdue boolean    not null default false
);

CREATE TABLE if not exists sub (
    id uuid primary key DEFAULT uuid_generate_v4(),
    name varchar(30) not null,
    todo_id uuid,
    foreign key (todo_id) references todo(id)
    ON UPDATE CASCADE ON DELETE cascade,
    completed boolean not null default false
);
