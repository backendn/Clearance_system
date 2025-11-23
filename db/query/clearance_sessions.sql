-- name: CreateSession :one
INSERT INTO clearance_sessions (
    name, start_date, end_date, active, created_at
) VALUES ($1,$2,$3,$4,NOW())
RETURNING *;

-- name: GetSession :one
SELECT * FROM clearance_sessions WHERE id = $1 LIMIT 1;

-- name: ListSessions :many
SELECT * FROM clearance_sessions ORDER BY id;

-- name: UpdateSession :one
UPDATE clearance_sessions SET
    name = $1,
    start_date = $2,
    end_date = $3,
    active = $4
WHERE id = $5
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM clearance_sessions WHERE id = $1;
