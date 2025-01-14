create table comments
(
    id         bigint                   not null primary key,
    article_id bigint                   not null references articles (id),
    user_id    bigint                   not null references users (id),
    reply_to   bigint                   null references comments (id),
    content    text                     not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);