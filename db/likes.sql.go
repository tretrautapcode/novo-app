// Code generated by sqlc. DO NOT EDIT.
// source: likes.sql

package db

import (
	"context"
)

const disLikes = `-- name: DisLikes :exec
INSERT INTO book_group_likes(user_id, book_group_id, point)
VALUES ($1, $2, -1)
`

type DisLikesParams struct {
	UserID      int32 `json:"userID"`
	BookGroupID int32 `json:"bookGroupID"`
}

func (q *Queries) DisLikes(ctx context.Context, arg DisLikesParams) error {
	_, err := q.db.Exec(ctx, disLikes, arg.UserID, arg.BookGroupID)
	return err
}

const getDislikes = `-- name: GetDislikes :one
SELECT coalesce(SUM(point), 0) as totalLikes FROM book_group_likes WHERE book_group_id = $1 AND point < 0
`

func (q *Queries) GetDislikes(ctx context.Context, bookGroupID int32) (interface{}, error) {
	row := q.db.QueryRow(ctx, getDislikes, bookGroupID)
	var coalesce interface{}
	err := row.Scan(&coalesce)
	return coalesce, err
}

const getLikes = `-- name: GetLikes :one
SELECT coalesce(SUM(point), 0) as totalLikes FROM book_group_likes WHERE book_group_id = $1 AND point > 0
`

func (q *Queries) GetLikes(ctx context.Context, bookGroupID int32) (interface{}, error) {
	row := q.db.QueryRow(ctx, getLikes, bookGroupID)
	var coalesce interface{}
	err := row.Scan(&coalesce)
	return coalesce, err
}

const likes = `-- name: Likes :exec
INSERT INTO book_group_likes(user_id, book_group_id, point)
VALUES ($1, $2, 1)
`

type LikesParams struct {
	UserID      int32 `json:"userID"`
	BookGroupID int32 `json:"bookGroupID"`
}

func (q *Queries) Likes(ctx context.Context, arg LikesParams) error {
	_, err := q.db.Exec(ctx, likes, arg.UserID, arg.BookGroupID)
	return err
}

const unlikes = `-- name: Unlikes :exec
DELETE FROM book_group_likes WHERE user_id = $1 AND book_group_id = $2
`

type UnlikesParams struct {
	UserID      int32 `json:"userID"`
	BookGroupID int32 `json:"bookGroupID"`
}

func (q *Queries) Unlikes(ctx context.Context, arg UnlikesParams) error {
	_, err := q.db.Exec(ctx, unlikes, arg.UserID, arg.BookGroupID)
	return err
}
