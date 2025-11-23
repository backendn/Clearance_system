-- name: CreateClearanceRecord :one
INSERT INTO clearance_records (
    student_id, clearance_item_id, session_id,
    status, note, handled_by,
    handled_at, attachment_url, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())
RETURNING *;

-- name: GetClearanceRecord :one
SELECT * FROM clearance_records WHERE id = $1 LIMIT 1;

-- name: ListRecordsByStudent :many
SELECT * FROM clearance_records
WHERE student_id = $1
ORDER BY clearance_item_id;

-- name: ListRecordsBySession :many
SELECT * FROM clearance_records
WHERE session_id = $1
ORDER BY student_id;

-- name: UpdateClearanceRecordStatus :one
UPDATE clearance_records SET
    status = $1,
    note = $2,
    handled_by = $3,
    handled_at = $4,
    attachment_url = $5,
    updated_at = NOW()
WHERE id = $6
RETURNING *;

-- name: DeleteClearanceRecord :exec
DELETE FROM clearance_records WHERE id = $1;
