alter table products
    add column tags text[] not null default '{}';