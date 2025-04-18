-- +goose Up
CREATE TABLE refresh_tokens as (
    token text PRIMARY key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null REFERENCES users(id) ON DELETE CASCADE,
    expires_at timestamp not null,
    revoked_at timestamp
);

-- +goose Down
DROP TABLE refresh_tokens;