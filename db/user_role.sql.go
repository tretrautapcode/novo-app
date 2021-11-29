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

const getRoleId = `-- name: GetRoleId :one
SELECT id FROM roles WHERE name = $1
`

func (q *Queries) GetRoleId(ctx context.Context, name string) (int32, error) {
	row := q.db.QueryRow(ctx, getRoleId, name)
	var id int32
	err := row.Scan(&id)
	return id, err
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

const role = `-- name: Role :one
SELECT r.name                             role_name,
       array_agg(module || '.' || action)::text[] role_permissions
FROM roles r
         LEFT JOIN role_permissions rp on r.id = rp.role_id
WHERE r.id = $1
GROUP BY r.name
`

type RoleRow struct {
	RoleName        string   `json:"roleName"`
	RolePermissions []string `json:"rolePermissions"`
}

func (q *Queries) Role(ctx context.Context, id int32) (RoleRow, error) {
	row := q.db.QueryRow(ctx, role, id)
	var i RoleRow
	err := row.Scan(&i.RoleName, &i.RolePermissions)
	return i, err
}
