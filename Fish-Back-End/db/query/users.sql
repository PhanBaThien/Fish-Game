-- name: CheckUserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: GetUserByUsername :one
SELECT id, username, email, password, role_id, created_at, updated_at
FROM users 
WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username, email, password, role_id, created_at, updated_at
FROM users 
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password, role_id)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;