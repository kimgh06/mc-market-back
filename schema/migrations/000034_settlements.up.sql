create table settlements(
    id bigint primary key generated always as identity,
    user_id bigint not null references users(id),
    amount numeric(10, 2) not null,
    month varchar not null default to_char(now(), 'YYYY-MM'),
    settled boolean not null default false
);