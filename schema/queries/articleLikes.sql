-- Insert a new like
insert into articles_likes (article_id, user_id, kind)
values ($1, $2, $3)
returning *;

-- Select a like by article_id and user_id
select * from articles_likes
where article_id = $1 and user_id = $2;

-- Update a like by article_id and user_id
update articles_likes
set kind = $3
where article_id = $1 and user_id = $2
returning *;

-- Delete a like by article_id and user_id
delete from articles_likes
where article_id = $1 and user_id = $2
returning *;
