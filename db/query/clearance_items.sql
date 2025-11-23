-- name: CreateClearanceItem :one
INSERT INTO clearance_items (
    code, title, description, department_id,
    approver_staff_id, requires_attachment,
    sequence, created_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
RETURNING *;

-- name: GetClearanceItem :one
SELECT * FROM clearance_items WHERE id = $1 LIMIT 1;

-- name: ListClearanceItems :many
SELECT * FROM clearance_items ORDER BY sequence;

-- name: ListItemsByDepartment :many
SELECT * FROM clearance_items WHERE department_id = $1 ORDER BY sequence;

-- name: UpdateClearanceItem :one
UPDATE clearance_items SET
    code = $1,
    title = $2,
    description = $3,
    department_id = $4,
    approver_staff_id = $5,
    requires_attachment = $6,
    sequence = $7
WHERE id = $8 RETURNING *;

-- name: DeleteClearanceItem :exec
DELETE FROM clearance_items WHERE id = $1;
