// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: roles.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createRole = `-- name: CreateRole :one
INSERT INTO roles (name, level) VALUES ($1, $2) RETURNING id, name, level, created_at
`

type CreateRoleParams struct {
	Name  string
	Level int32
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, createRole, arg.Name, arg.Level)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Level,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAllRoles = `-- name: DeleteAllRoles :exec
DELETE FROM roles
`

func (q *Queries) DeleteAllRoles(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteAllRoles)
	return err
}

const getAllRoles = `-- name: GetAllRoles :many
SELECT id, name, level, created_at FROM roles
`

func (q *Queries) GetAllRoles(ctx context.Context) ([]Role, error) {
	rows, err := q.db.Query(ctx, getAllRoles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Level,
			&i.CreatedAt,
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

const getRoleById = `-- name: GetRoleById :one
SELECT id, name, level, created_at FROM roles WHERE id = $1
`

func (q *Queries) GetRoleById(ctx context.Context, id uuid.UUID) (Role, error) {
	row := q.db.QueryRow(ctx, getRoleById, id)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Level,
		&i.CreatedAt,
	)
	return i, err
}

const getRoleByName = `-- name: GetRoleByName :one
SELECT id, name, level, created_at FROM roles WHERE name = $1
`

func (q *Queries) GetRoleByName(ctx context.Context, name string) (Role, error) {
	row := q.db.QueryRow(ctx, getRoleByName, name)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Level,
		&i.CreatedAt,
	)
	return i, err
}