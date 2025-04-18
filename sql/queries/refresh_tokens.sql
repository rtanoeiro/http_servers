-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, NOW(), NOW(), $2, $3, null)

RETURNING *;

-- name: GetRefreshToken :one
select
    token,
    user_id
from refresh_tokens
where token = $1
and revoked_at is null;

-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = NOW(), updated_at = NOW()
where token = $1;