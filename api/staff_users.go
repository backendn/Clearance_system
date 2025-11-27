package api

import (
	"net/http"

	sqlc "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

// ----------------------
// REQUEST MODELS
// ----------------------

type createStaffUserRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	Email        string `json:"email" binding:"required,email"`
	FullName     string `json:"full_name" binding:"required"`
	DepartmentID int64  `json:"department_id" binding:"required,min=1"`
	RoleID       int64  `json:"role_id" binding:"required,min=1"`
	Password     string `json:"password" binding:"required,min=6"`
}

type updateStaffUserRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	Email        string `json:"email" binding:"required,email"`
	FullName     string `json:"full_name" binding:"required"`
	DepartmentID int64  `json:"department_id" binding:"required,min=1"`
	RoleID       int64  `json:"role_id" binding:"required,min=1"`
	Password     string `json:"password" binding:"required,min=6"`
}

// ----------------------
// RESPONSE MODEL
// ----------------------

type staffUserResponse struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	DepartmentID int64  `json:"department_id"`
	RoleID       int64  `json:"role_id"`
	CreatedAt    string `json:"created_at"`
}

func newStaffUserResponse(u sqlc.StaffUser) staffUserResponse {
	return staffUserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		FullName:     u.FullName,
		DepartmentID: u.DepartmentID,
		RoleID:       u.RoleID,
		CreatedAt:    u.CreatedAt.String(),
	}
}
func (server *Server) CreateStaffUser(ctx *gin.Context) {
	var req createStaffUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check email exists
	_, err := server.store.GetStaffUserByEmail(ctx, req.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("email already exists"))
		return
	}

	// Check username exists
	_, err = server.store.GetStaffUserByUsername(ctx, req.Username)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("username already exists"))
		return
	}

	// Check department exists
	_, err = server.store.GetDepartment(ctx, req.DepartmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("department_id does not exist"))
		return
	}

	// Hash password
	hashed, err := HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.CreateStaffUserParams{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		DepartmentID: req.DepartmentID,
		RoleID:       req.RoleID,
		PasswordHash: hashed,
	}

	user, err := server.store.CreateStaffUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, newStaffUserResponse(user))
}
func (server *Server) GetStaffUser(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}

	user, err := server.store.GetStaffUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorMessage("staff user not found"))
		return
	}

	ctx.JSON(http.StatusOK, newStaffUserResponse(user))
}
func (server *Server) ListStaffUsers(ctx *gin.Context) {
	limit, offset := getPagination(ctx)

	users, err := server.store.ListStaffUsers(ctx, sqlc.ListStaffUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := make([]staffUserResponse, len(users))
	for i, u := range users {
		resp[i] = newStaffUserResponse(u)
	}

	ctx.JSON(http.StatusOK, resp)
}
func (server *Server) UpdateStaffUser(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateStaffUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify department exists
	_, err = server.store.GetDepartment(ctx, req.DepartmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("department_id does not exist"))
		return
	}

	hashed, err := HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.UpdateStaffUserParams{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		DepartmentID: req.DepartmentID,
		RoleID:       req.RoleID,
		PasswordHash: hashed,
		ID:           id,
	}

	user, err := server.store.UpdateStaffUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newStaffUserResponse(user))
}
func (server *Server) DeleteStaffUser(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteStaffUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorMessage("staff user not found"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "staff user deleted"})
}
