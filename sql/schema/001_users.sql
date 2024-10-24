-- +goose Up
CREATE TABLE users(
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name    TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;