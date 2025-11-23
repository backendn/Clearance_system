-- name: CreateDepartment :one
INSERT INTO departments (
    code, name, created_at
) VALUES ($1,$2,NOW())
RETURNING *;

-- name: GetDepartment :one
SELECT * FROM departments WHERE id = $1 LIMIT 1;

-- name: GetDepartmentByCode :one
SELECT * FROM departments WHERE code = $1 LIMIT 1;

-- name: ListDepartments :many
SELECT * FROM departments ORDER BY id;

-- name: UpdateDepartment :one
UPDATE departments SET
    code = $1,
    name = $2
WHERE id = $3 RETURNING *;

-- name: DeleteDepartment :exec
DELETE FROM departments WHERE id = $1;
