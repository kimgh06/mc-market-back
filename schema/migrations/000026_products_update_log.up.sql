create table products_update_log (
    id serial primary key,
    product_id bigint not null,
    title varchar(100) not null,
    content text not null,
    updated_at timestamp not null default now()
);