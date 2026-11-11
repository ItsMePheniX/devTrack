-- name: CreateUser :one
INSERT INTO users(
    email,
    password_hash,
    full_name,
    role
)VALUES(
    $1,$2,$3,$4
)
RETURNING id,email,full_name,role,avatar_url,created_at,updated_at;

-- name: GetUserByEmail :one
SELECT id,email,password_hash,full_name,role,avatar_url,created_at,updated_at
FROM users
where email=$1 LIMIT 1;

-- name: GetUserByID :one
SELECT id,email,full_name,role,avatar_url,created_at,updated_at
FROM users
WHERE id=$1 LIMIT 1;