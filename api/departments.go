package api

import (
	"net/http"

	sqlc "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

//
// ===============================
//  REQUEST MODELS
// ===============================
//

type CreateDepartmentRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type UpdateDepartmentRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

//
// ===============================
//  RESPONSE FORMATTER
// ===============================
//

func newDepartmentResponse(d sqlc.Department) gin.H {
	return gin.H{
		"id":         d.ID,
		"code":       d.Code,
		"name":       d.Name,
		"created_at": d.CreatedAt,
	}
}

//
// ===============================
//  HANDLERS
// ===============================
//

// ---------- CREATE DEPARTMENT ----------
func (server *Server) CreateDepartment(ctx *gin.Context) {
	var req CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Ensure unique code
	_, err := server.store.GetDepartmentByCode(ctx, req.Code)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("department code already exists"))
		return
	}

	arg := sqlc.CreateDepartmentParams{
		Code: req.Code,
		Name: req.Name,
	}

	dept, err := server.store.CreateDepartment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, newDepartmentResponse(dept))
}

// ---------- GET DEPARTMENT ----------
func (server *Server) GetDepartment(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dept, err := server.store.GetDepartment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorMessage("department not found"))
		return
	}

	ctx.JSON(http.StatusOK, newDepartmentResponse(dept))
}

// ---------- LIST DEPARTMENTS ----------
func (server *Server) ListDepartments(ctx *gin.Context) {
	depts, err := server.store.ListDepartments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	list := make([]gin.H, len(depts))
	for i, d := range depts {
		list[i] = newDepartmentResponse(d)
	}

	ctx.JSON(http.StatusOK, list)
}

// ---------- UPDATE DEPARTMENT ----------
func (server *Server) UpdateDepartment(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.UpdateDepartmentParams{
		ID:   id,
		Code: req.Code,
		Name: req.Name,
	}

	dept, err := server.store.UpdateDepartment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newDepartmentResponse(dept))
}

// ---------- DELETE DEPARTMENT ----------
func (server *Server) DeleteDepartment(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteDepartment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
