// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const bookGroupsByUser = `-- name: BookGroupsByUser :many
SELECT bg.id,
       (array_agg(i.path))[1]   AS image,
       (array_agg(bg.title))[1] as title,
       bct.latest_chapter,
       bct.last_updated,
       bct.views,
       bcm.comments,
       bgl.likes
FROM book_groups as bg
         JOIN users u on u.id = bg.owner_id
         LEFT JOIN Lateral (
    SELECT count(bcm.id) AS comments
    FROM book_comments bcm
    WHERE bcm.book_group_id = bg.id
    ) bcm ON TRUE
         LEFT JOIN LATERAL (
    SELECT coalesce(sum(bgl.point), 0) AS likes
    FROM book_group_likes bgl
    WHERE bgl.book_group_id = bg.id
    ) bgl ON TRUE
         LEFT JOIN LATERAL (
    SELECT (array_agg(bct.chapter_number ORDER BY bct.date_created DESC))[1] AS latest_chapter,
           MAX(bct.date_created)                                             AS last_updated,
           coalesce(sum(bcv.count), 0)                                       AS views
    FROM book_chapters bct
             LEFT JOIN book_chapter_views bcv
                       ON bct.id = bcv.book_chapter_id
    WHERE bct.book_group_id = bg.id
    ) bct ON TRUE
         LEFT JOIN images i on bg.primary_cover_art_id = i.id
WHERE u.id = $1
GROUP BY bg.id, bg.title, i.path, bct.latest_chapter, bct.last_updated, bct.views, bcm.comments, bgl.likes
ORDER BY last_updated DESC NULLS LAST
`

type BookGroupsByUserRow struct {
	ID            int32       `json:"id"`
	Image         interface{} `json:"image"`
	Title         interface{} `json:"title"`
	LatestChapter interface{} `json:"latestChapter"`
	LastUpdated   interface{} `json:"lastUpdated"`
	Views         interface{} `json:"views"`
	Comments      int64       `json:"comments"`
	Likes         interface{} `json:"likes"`
}

func (q *Queries) BookGroupsByUser(ctx context.Context, id int32) ([]BookGroupsByUserRow, error) {
	rows, err := q.db.Query(ctx, bookGroupsByUser, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BookGroupsByUserRow
	for rows.Next() {
		var i BookGroupsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Image,
			&i.Title,
			&i.LatestChapter,
			&i.LastUpdated,
			&i.Views,
			&i.Comments,
			&i.Likes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const checkEmailExist = `-- name: CheckEmailExist :one
SELECT exists(select 1 from users where email = $1)
`

func (q *Queries) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, checkEmailExist, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUsernameExist = `-- name: CheckUsernameExist :one
SELECT exists(select 1 from users where user_name = $1)
`

func (q *Queries) CheckUsernameExist(ctx context.Context, userName sql.NullString) (bool, error) {
	row := q.db.QueryRow(ctx, checkUsernameExist, userName)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const completeOauthAccount = `-- name: CompleteOauthAccount :exec
UPDATE users
SET user_name       = $2,
    avatar_image_id = $3,
    role_id         = $4
WHERE id = $1
`

type CompleteOauthAccountParams struct {
	ID            int32          `json:"id"`
	UserName      sql.NullString `json:"userName"`
	AvatarImageID sql.NullInt32  `json:"avatarImageID"`
	RoleID        int32          `json:"roleID"`
}

func (q *Queries) CompleteOauthAccount(ctx context.Context, arg CompleteOauthAccountParams) error {
	_, err := q.db.Exec(ctx, completeOauthAccount,
		arg.ID,
		arg.UserName,
		arg.AvatarImageID,
		arg.RoleID,
	)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE user_name = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userName sql.NullString) error {
	_, err := q.db.Exec(ctx, deleteUser, userName)
	return err
}

const getUserInfo = `-- name: GetUserInfo :one
SELECT users.id,
       users.user_name,
       users.email,
       users.password,
       r.name as role,
       users.summary,
       i.path as avatarPath
FROM users
         JOIN roles r on users.role_id = r.id
         LEFT JOIN images i on users.avatar_image_id = i.id
WHERE users.id = $1
`

type GetUserInfoRow struct {
	ID         int32          `json:"id"`
	UserName   sql.NullString `json:"userName"`
	Email      string         `json:"email"`
	Password   sql.NullString `json:"password"`
	Role       string         `json:"role"`
	Summary    sql.NullString `json:"summary"`
	Avatarpath sql.NullString `json:"avatarpath"`
}

func (q *Queries) GetUserInfo(ctx context.Context, id int32) (GetUserInfoRow, error) {
	row := q.db.QueryRow(ctx, getUserInfo, id)
	var i GetUserInfoRow
	err := row.Scan(
		&i.ID,
		&i.UserName,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.Summary,
		&i.Avatarpath,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO users(user_name, password, email, role_id)
VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
RETURNING id, date_created, user_name, password, email, summary, avatar_image_id, role_id, favorite_list
`

type InsertUserParams struct {
	UserName sql.NullString `json:"userName"`
	Password sql.NullString `json:"password"`
	Email    string         `json:"email"`
	RoleName string         `json:"roleName"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (User, error) {
	row := q.db.QueryRow(ctx, insertUser,
		arg.UserName,
		arg.Password,
		arg.Email,
		arg.RoleName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DateCreated,
		&i.UserName,
		&i.Password,
		&i.Email,
		&i.Summary,
		&i.AvatarImageID,
		&i.RoleID,
		&i.FavoriteList,
	)
	return i, err
}

const searchUsers = `-- name: SearchUsers :many
SELECT users.user_name, users.id, i.path
FROM users
         LEFT JOIN images i on users.avatar_image_id = i.id
WHERE user_name ILIKE '%' || $1 || '%'
LIMIT 5
`

type SearchUsersRow struct {
	UserName sql.NullString `json:"userName"`
	ID       int32          `json:"id"`
	Path     sql.NullString `json:"path"`
}

func (q *Queries) SearchUsers(ctx context.Context, dollar_1 sql.NullString) ([]SearchUsersRow, error) {
	rows, err := q.db.Query(ctx, searchUsers, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchUsersRow
	for rows.Next() {
		var i SearchUsersRow
		if err := rows.Scan(&i.UserName, &i.ID, &i.Path); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users
SET password = $2
WHERE id = $1
`

type UpdatePasswordParams struct {
	ID       int32          `json:"id"`
	Password sql.NullString `json:"password"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, updatePassword, arg.ID, arg.Password)
	return err
}

const updateUserInfo = `-- name: UpdateUserInfo :exec
Update users
SET email     = $2,
    user_name = $3,
    summary   = $4
WHERE id = $1
`

type UpdateUserInfoParams struct {
	ID       int32          `json:"id"`
	Email    string         `json:"email"`
	UserName sql.NullString `json:"userName"`
	Summary  sql.NullString `json:"summary"`
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) error {
	_, err := q.db.Exec(ctx, updateUserInfo,
		arg.ID,
		arg.Email,
		arg.UserName,
		arg.Summary,
	)
	return err
}

const userByEmail = `-- name: UserByEmail :one
SELECT id, date_created, user_name, password, email, summary, avatar_image_id, role_id, favorite_list
FROM users
WHERE email = $1
    FETCH FIRST ROWS ONLY
`

func (q *Queries) UserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, userByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DateCreated,
		&i.UserName,
		&i.Password,
		&i.Email,
		&i.Summary,
		&i.AvatarImageID,
		&i.RoleID,
		&i.FavoriteList,
	)
	return i, err
}

const userByUsernameOrEmail = `-- name: UserByUsernameOrEmail :one
SELECT id, date_created, user_name, password, email, summary, avatar_image_id, role_id, favorite_list
FROM users
WHERE user_name = $1
   OR email = $1
    FETCH FIRST ROWS ONLY
`

func (q *Queries) UserByUsernameOrEmail(ctx context.Context, userName sql.NullString) (User, error) {
	row := q.db.QueryRow(ctx, userByUsernameOrEmail, userName)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DateCreated,
		&i.UserName,
		&i.Password,
		&i.Email,
		&i.Summary,
		&i.AvatarImageID,
		&i.RoleID,
		&i.FavoriteList,
	)
	return i, err
}
