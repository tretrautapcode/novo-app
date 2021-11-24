// Code generated by sqlc. DO NOT EDIT.
// source: user_role.sql

package db

import (
	"context"
	"database/sql"
)

const deleteRole = `-- name: DeleteRole :exec
DELETE
FROM roles
WHERE name = $1
`

func (q *Queries) DeleteRole(ctx context.Context, name string) error {
	_, err := q.db.Exec(ctx, deleteRole, name)
	return err
}

const insertNewRole = `-- name: InsertNewRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING id, name, description
`

type InsertNewRoleParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) InsertNewRole(ctx context.Context, arg InsertNewRoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, insertNewRole, arg.Name, arg.Description)
	var i Role
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const roleByUserId = `-- name: RoleByUserId :one
SELECT r.name                             role_name,
       array_agg(module || '.' || action) role_permissions
FROM role_permissions rp
         JOIN roles r on r.id = rp.role_id
         JOIN users u on r.id = u.role_id
WHERE u.id = $1
GROUP BY r.name
`

type RoleByUserIdRow struct {
	RoleName        string      `json:"roleName"`
	RolePermissions interface{} `json:"rolePermissions"`
}

func (q *Queries) RoleByUserId(ctx context.Context, id int32) (RoleByUserIdRow, error) {
	row := q.db.QueryRow(ctx, roleByUserId, id)
	var i RoleByUserIdRow
	err := row.Scan(&i.RoleName, &i.RolePermissions)
	return i, err
}
