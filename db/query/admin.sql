-- name: CreateAdmin :one
INSERT INTO admins (
    username,
    hashed_password,
    full_name,
    email,
    role,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAdmin :one
SELECT *
FROM admins
WHERE id = $1
LIMIT 1;

-- name: GetAdminByUsername :one
SELECT *
FROM admins
WHERE username = $1
LIMIT 1;

-- name: ListAdmins :many
SELECT *
FROM admins
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateAdmin :one
UPDATE admins
SET
    full_name = $2,
    email = $3,
    role = $4,
    is_active = $5
WHERE id = $1
RETURNING *;

-- name: UpdateAdminPassword :exec
UPDATE admins
SET hashed_password = $2
WHERE id = $1;
-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE id = $1;

-- name: AdminExistsByUsername :one
SELECT EXISTS(
    SELECT 1 FROM admins WHERE username = $1
) AS exists;
