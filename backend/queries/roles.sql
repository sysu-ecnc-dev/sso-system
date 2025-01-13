-- name: GetAllRoles :many
SELECT * FROM roles;

-- name: CreateRole :one
INSERT INTO roles (name, level) VALUES ($1, $2) RETURNING *;

-- name: DeleteAllRoles :exec
DELETE FROM roles;

-- name: GetRoleById :one
SELECT * FROM roles WHERE id = $1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;