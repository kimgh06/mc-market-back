-- Insert a new comment
insert into comments (article_id, user_id, reply_to, content)
values ($1, $2, $3, $4)
returning *;

-- Select a comment by id
select * from comments
where id = $1;

-- Update a comment by id
update comments
set content = $2, updated_at = now()
where id = $1
returning *;

-- Delete a comment by id
delete from comments
where id = $1
returning *;