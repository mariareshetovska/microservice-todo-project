-- migrate up
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

-- migrate down
DROP TABLE todo;