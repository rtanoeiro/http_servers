-- name: InsertChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)

RETURNING *;

-- name: GetAllChirps :many
SELECT 
    id,
    created_at,
    updated_at,
    body
FROM chirps
ORDER BY created_at;

-- name: GetSingleChirp :one
SELECT 
    id,
    created_at,
    updated_at,
    body,
    user_id
FROM chirps
WHERE id = $1;

-- name: GetAuthorChirps :many
SELECT
    id,
    created_at,
    updated_at,
    body
FROM chirps
WHERE user_id = $1
ORDER by created_at;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;