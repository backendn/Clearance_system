-- name: CreateClearanceRequest :one
INSERT INTO clearance_requests (
    student_id, session_id
) VALUES ($1, $2)
RETURNING *;

-- name: GetClearanceRequest :one
SELECT * FROM clearance_requests WHERE id = $1 LIMIT 1;

-- name: GetStudentRequestForSession :one
SELECT * FROM clearance_requests
WHERE student_id = $1 AND session_id = $2
LIMIT 1;

-- name: ListRequestsByStudent :many
SELECT * FROM clearance_requests
WHERE student_id = $1
ORDER BY created_at DESC;

-- name: UpdateClearanceRequestStatus :exec
UPDATE clearance_requests
SET status = $1
WHERE id = $2;

-- name: ListAllRequests :many
SELECT * FROM clearance_requests
ORDER BY created_at DESC;
