-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id uuid REFERENCES users(id) not null
);

-- +goose Down
DROP TABLE chirps;