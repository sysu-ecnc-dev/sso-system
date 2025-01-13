-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: CreateUser :one
INSERT INTO
    users (
        username,
        password_hash,
        full_name,
        email,
        role_id
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: DeleteUserByUsername :exec
DELETE FROM users WHERE username = $1;