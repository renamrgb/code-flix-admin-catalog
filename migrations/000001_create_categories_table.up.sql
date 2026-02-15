create table categories (
    id varchar(36) not null primary key,
    name varchar(255) not null,
    description varchar(4000),
    activated boolean not null default true,
    created_at datetime(6) not null,
    updated_at datetime(6) not null,
    deleted_at datetime(6)
);
