alter table products_versions
    add column updated_at timestamp not null default now();