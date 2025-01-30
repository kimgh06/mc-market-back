CREATE TABLE article_head_type
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

alter table articles
    add constraint fk_article_head
        foreign key (head) references article_head_type (id);