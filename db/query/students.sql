-- name: CreateStudent :one
INSERT INTO students (
    student_number,
    first_name,
    last_name,
    email,
    phone,
    department_id,
    enrollment_year,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
)
RETURNING *;

-- name: GetStudent :one
SELECT * FROM students
WHERE id = $1 LIMIT 1;

-- name: GetStudentByStudentNumber :one
SELECT * FROM students
WHERE student_number = $1 LIMIT 1;

-- name: ListStudents :many
SELECT * FROM students
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateStudent :one
UPDATE students SET
    student_number = $1,
    first_name = $2,
    last_name = $3,
    email = $4,
    phone = $5,
    department_id = $6,
    enrollment_year = $7
WHERE id = $8
RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE id = $1;
