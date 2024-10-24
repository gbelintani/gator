-- +goose Up
CREATE TABLE posts(
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at timestamp NOT NULL,
    feed_id uuid NOT NULL,
    CONSTRAINT fk_feed FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE posts;