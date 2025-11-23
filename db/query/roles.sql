-- name: CreateRole :one
INSERT INTO roles (name)
VALUES ($1)
RETURNING *;

-- name: GetRole :one
SELECT * FROM roles WHERE id = $1 LIMIT 1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY id;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;
