create table articles_likes (
    article_id bigint not null references articles (id),
    user_id bigint not null references users (id),
    kind boolean not null,
    created_at timestamp with time zone not null default now(),
    primary key (article_id, user_id)
);