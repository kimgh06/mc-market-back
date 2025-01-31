DROP TABLE IF EXISTS article_head_type;

alter table articles
drop
constraint fk_article_head;