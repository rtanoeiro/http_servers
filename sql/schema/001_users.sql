-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text unique not null,
    hashed_password text not null,
    is_chirpy_red bool not null
);

-- +goose Down
DROP TABLE users;