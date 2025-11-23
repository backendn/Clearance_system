-- name: CreateStaffUser :one
INSERT INTO staff_users (
    username, email, full_name, department_id,
    role_id, password_hash, created_at
) VALUES ($1,$2,$3,$4,$5,$6,NOW())
RETURNING *;

-- name: GetStaffUser :one
SELECT * FROM staff_users WHERE id = $1 LIMIT 1;

-- name: GetStaffUserByUsername :one
SELECT * FROM staff_users WHERE username = $1 LIMIT 1;

-- name: GetStaffUserByEmail :one
SELECT * FROM staff_users WHERE email = $1 LIMIT 1;

-- name: ListStaffUsers :many
SELECT * FROM staff_users ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateStaffUser :one
UPDATE staff_users SET
    username = $1,
    email = $2,
    full_name = $3,
    department_id = $4,
    role_id = $5,
    password_hash = $6
WHERE id = $7
RETURNING *;

-- name: DeleteStaffUser :exec
DELETE FROM staff_users WHERE id = $1;
