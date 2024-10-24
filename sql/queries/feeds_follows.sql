-- name: CreateFollow :one
WITH inserted AS(
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
) 

SELECT i.*, u.name as user_name, f.name as feed_name from inserted i
join users u on u.id = i.user_id
join feeds f on f.id = i.feed_id;


-- name: GetFeedFollowsForUser :many
SELECT fw.*, u.name as user_name, f.name as feed_name FROM feed_follows fw
JOIN users u on u.id = fw.user_id
JOIN feeds f on f.id = fw.feed_id
WHERE fw.user_id = $1;