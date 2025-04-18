-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    $1, $2, $3, $4, $5, false
)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: CheckUserWithEmail :one
SELECT
    id,
    hashed_password,
    created_at,
    updated_at,
    is_chirpy_red
from users
where email = $1;

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = NOW(),
    email = $1,
    hashed_password = $2
WHERE id = $3

RETURNING *;

-- name: CheckUserWithID :one
SELECT
    id,
    is_chirpy_red
from users
where id = $1;

-- name: UpgradeToChirpy :one

UPDATE users
set is_chirpy_red = true
where id = $1

RETURNING id;