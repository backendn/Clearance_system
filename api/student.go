package api

import (
	"net/http"
	"strconv"

	sqlc "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateStudentRequest is the expected body for creating a student
type CreateStudentRequest struct {
	StudentNumber  string `json:"student_number" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required"`
	DepartmentID   int64  `json:"department_id" binding:"required"`
	EnrollmentYear int32  `json:"enrollment_year" binding:"required"`
}

// CreateStudent handler
func (server *Server) CreateStudent(ctx *gin.Context) {
	var req CreateStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := sqlc.CreateStudentParams{
		StudentNumber:  req.StudentNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DepartmentID:   req.DepartmentID,
		EnrollmentYear: req.EnrollmentYear,
	}

	// 1️⃣ Validate department exists
	_, err := server.store.GetDepartment(ctx, req.DepartmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid department_id: department does not exist",
		})
		return
	}

	student, err := server.store.CreateStudent(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// GetStudent handler
func (server *Server) GetStudentByNumber(ctx *gin.Context) {
	studentNumber := ctx.Param("student_number")
	if studentNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorMessage("missing student_number"))
		return
	}

	student, err := server.store.GetStudentByStudentNumber(ctx, studentNumber)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorMessage("student not found"))
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// UpdateStudent handler
func (server *Server) UpdateStudent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req CreateStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := sqlc.UpdateStudentParams{
		StudentNumber:  req.StudentNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DepartmentID:   req.DepartmentID,
		EnrollmentYear: req.EnrollmentYear,
		ID:             id,
	}

	student, err := server.store.UpdateStudent(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// DeleteStudent handler
func (server *Server) DeleteStudent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = server.store.DeleteStudent(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "student deleted"})
}
