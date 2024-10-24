-- +goose Up
ALTER TABLE feeds
ADD COLUMN fetched_at timestamp NULL;


-- +goose Down
ALTER TABLE feeds
DROP COLUMN fetched_at;