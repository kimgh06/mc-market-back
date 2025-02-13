create table products_versions (
    id serial primary key,
    product_id bigint not null,
    version_name varchar(100) not null,
    link varchar(100) not null
);