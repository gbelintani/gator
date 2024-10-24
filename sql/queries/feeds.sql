-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT f.*, u.name as user_name FROM feeds f
JOIN users u on u.id = f.user_id;

-- name: GetByUrl :one
select * from feeds 
where url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from feeds 
ORDER BY fetched_at asc NULLS FIRST
LIMIT 1;