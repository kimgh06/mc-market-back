alter table articles
    add column author bigint not null default 0 references users (id);