create table payments
(
    id         bigint primary key,
    agent      bigint                   not null references users (id),
    order_id   uuid                     not null unique default gen_random_uuid(),
    amount     int                      not null check ( amount <> 0 ),
    approved   bool                     not null        default false,
    created_at timestamp with time zone not null        default now()
);