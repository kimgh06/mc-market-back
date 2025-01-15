alter table articles
    add column head integer null default null,
    add constraint fk_article_head
    foreign key (head) references article_head_type(id);