alter table articles
    add column comment_disabled boolean not null default false,
    add column like_disabled boolean not null default false;