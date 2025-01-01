alter table products
    add column details text not null default ''
        constraint details_length check ( char_length(details) <= 500 );