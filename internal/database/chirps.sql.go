// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chirps.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const deleteChirp = `-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1
`

func (q *Queries) DeleteChirp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteChirp, id)
	return err
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT 
    id,
    created_at,
    updated_at,
    body
FROM chirps
ORDER BY created_at
`

type GetAllChirpsRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
}

func (q *Queries) GetAllChirps(ctx context.Context) ([]GetAllChirpsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllChirpsRow
	for rows.Next() {
		var i GetAllChirpsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAuthorChirps = `-- name: GetAuthorChirps :many
SELECT
    id,
    created_at,
    updated_at,
    body
FROM chirps
WHERE user_id = $1
ORDER by created_at
`

type GetAuthorChirpsRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
}

func (q *Queries) GetAuthorChirps(ctx context.Context, userID uuid.UUID) ([]GetAuthorChirpsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAuthorChirps, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAuthorChirpsRow
	for rows.Next() {
		var i GetAuthorChirpsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSingleChirp = `-- name: GetSingleChirp :one
SELECT 
    id,
    created_at,
    updated_at,
    body,
    user_id
FROM chirps
WHERE id = $1
`

func (q *Queries) GetSingleChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getSingleChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const insertChirp = `-- name: InsertChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)

RETURNING id, created_at, updated_at, body, user_id
`

type InsertChirpParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	UserID    uuid.UUID
}

func (q *Queries) InsertChirp(ctx context.Context, arg InsertChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, insertChirp,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Body,
		arg.UserID,
	)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}
