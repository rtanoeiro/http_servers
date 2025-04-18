-- name: CreateRefreshToken
INSERT INTO refresh_tokens values (
    $1, NOW(), NOW(), $2, $3, null
)