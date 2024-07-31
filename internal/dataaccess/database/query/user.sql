-- name: CreateUser :one
INSERT INTO users (
    user_name,
    full_name,
    mail,
    role
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE mail = $1 LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateGmailUser :one
UPDATE users
SET mail = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
