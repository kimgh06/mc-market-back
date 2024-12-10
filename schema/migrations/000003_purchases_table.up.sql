create table purchases(
    id bigint not null,

    purchaser bigint not null references users(id),
    product bigint not null references products(id),

    purchased_at timestamp with time zone not null default now()
);