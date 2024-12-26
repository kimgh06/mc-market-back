create table articles
(
    id         bigint                   not null primary key,
    title      text                     not null,
    content    text                     not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);