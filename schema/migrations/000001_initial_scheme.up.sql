create
    extension pgcrypto;

create table public.users
(
    id          bigint primary key,
    nickname    varchar(64)              null unique,

    permissions int                      not null default 0,

    created_at  timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now()
);

create table public.products
(
    id             bigint primary key,
    creator        bigint                   not null references users (id),

    category       text                     not null,

    name           text                     not null
        constraint title_check check ( char_length(name) <= 50 ),
    description    text                     not null,
    usage          text                     not null,

    price          int                      not null,
    price_discount int                      null,

    created_at     timestamp with time zone not null default now(),
    updated_at     timestamp with time zone not null default now()
);

create table public.likes
(
    user_id    bigint references users (id),
    product_id bigint references products (id),

    unique (user_id, product_id)
);

create table public.downloads
(
    user_id    bigint references users (id)    not null,
    product_id bigint references products (id) not null,

    time       timestamp with time zone        not null default now(),
    unique (user_id, product_id)
);