-- name: CreateNotification :one
INSERT INTO notifications (
    recipient_user_id, recipient_student_id,
    message, read, created_at
) VALUES ($1,$2,$3,$4,NOW())
RETURNING *;

-- name: GetNotification :one
SELECT * FROM notifications WHERE id = $1 LIMIT 1;

-- name: ListNotificationsForUser :many
SELECT *
FROM notifications
WHERE recipient_user_id = $1
ORDER BY created_at DESC;

-- name: ListNotificationsForStudent :many
SELECT *
FROM notifications
WHERE recipient_student_id = $1
ORDER BY created_at DESC;

-- name: MarkNotificationRead :one
UPDATE notifications SET read = TRUE
WHERE id = $1 RETURNING *;

-- name: DeleteNotification :exec
DELETE FROM notifications WHERE id = $1;
