-- name: CreateFollow :one
WITH inserted AS(
  INSERT INTO feeds (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;
) 

SELECT inserted.*, u.name as user_name, f.name as feed_name
join users u on u.id = inserted.user_id
join feeds f on f.id = inserted.feed_id;
