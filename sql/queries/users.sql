-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: CheckUserLogin :one
SELECT
    id,
    hashed_password,
    created_at,
    updated_at
from users
where email = $1;