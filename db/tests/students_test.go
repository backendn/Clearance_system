package tests

import (
	"context"
	"testing"
	"time"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/backendn/clearance_system/util"
	"github.com/stretchr/testify/require"
)

// helper: create random department for FK
func createRandomDepartment(t *testing.T) db.Department {
	arg := db.CreateDepartmentParams{
		Code: util.RandomString(5),
		Name: util.RandomString(10),
	}

	dept, err := testQueries.CreateDepartment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dept)

	return dept
}

// helper: create random student
func createRandomStudent(t *testing.T) db.Student {
	dept := createRandomDepartment(t)

	arg := db.CreateStudentParams{
		StudentNumber:  util.RandomStudentNumber(),
		FirstName:      util.RandomString(6),
		LastName:       util.RandomString(6),
		Email:          util.RandomEmail(),
		Phone:          util.RandomPhone(),
		DepartmentID:   dept.ID,
		EnrollmentYear: 2021,
	}

	student, err := testQueries.CreateStudent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, student)

	return student
}

// -------------------
//        TESTS
// -------------------

func TestCreateStudent(t *testing.T) {
	student := createRandomStudent(t)
	require.NotEmpty(t, student)
}

func TestGetStudent(t *testing.T) {
	student1 := createRandomStudent(t)
	student2, err := testQueries.GetStudent(context.Background(), student1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, student2)
	require.Equal(t, student1.ID, student2.ID)
	require.Equal(t, student1.StudentNumber, student2.StudentNumber)
	require.Equal(t, student1.Email, student2.Email)
	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
}

func TestUpdateStudentPartial(t *testing.T) {
	student := createRandomStudent(t)

	newEmail := util.RandomEmail()
	newPhone := util.RandomPhone()

	arg := db.UpdateStudentPartialParams{
		ID:    student.ID,
		Email: newEmail,
		Phone: newPhone,
	}

	updated, err := testQueries.UpdateStudentPartial(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, student.ID, updated.ID)
	require.Equal(t, newEmail, updated.Email)
	require.Equal(t, newPhone, updated.Phone)
	// other fields remain unchanged
	require.Equal(t, student.StudentNumber, updated.StudentNumber)
	require.Equal(t, student.DepartmentID, updated.DepartmentID)
}

func TestDeleteStudent(t *testing.T) {
	student := createRandomStudent(t)

	err := testQueries.DeleteStudent(context.Background(), student.ID)
	require.NoError(t, err)

	// confirm deletion
	student2, err := testQueries.GetStudent(context.Background(), student.ID)
	require.Error(t, err)
	require.Empty(t, student2)
}
