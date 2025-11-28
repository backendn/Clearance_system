package api

import (
	"net/http"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createClearanceItemRequest struct {
	Code               string `json:"code" binding:"required"`
	Title              string `json:"title" binding:"required"`
	Description        string `json:"description"`
	DepartmentID       int64  `json:"department_id" binding:"required,min=1"`
	ApproverStaffID    int64  `json:"approver_staff_id" binding:"required,min=1"`
	RequiresAttachment bool   `json:"requires_attachment"`
	Sequence           int32  `json:"sequence" binding:"required,min=1"`
}

// POST /clearance-items
func (s *Server) createClearanceItem(ctx *gin.Context) {
	var req createClearanceItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate department exists
	_, err := s.store.GetDepartment(ctx, req.DepartmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid department ID"))
		return
	}

	// Validate approver staff exists
	_, err = s.store.GetStaffUser(ctx, req.ApproverStaffID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid approver staff ID"))
		return
	}

	arg := db.CreateClearanceItemParams{
		Code:               req.Code,
		Title:              req.Title,
		Description:        req.Description,
		DepartmentID:       req.DepartmentID,
		ApproverStaffID:    req.ApproverStaffID,
		RequiresAttachment: req.RequiresAttachment,
		Sequence:           req.Sequence,
	}

	item, err := s.store.CreateClearanceItem(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

// GET /clearance-items/:id
func (s *Server) getClearanceItem(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	item, err := s.store.GetClearanceItem(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorMessage("clearance item not found"))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// GET /clearance-items
func (s *Server) listClearanceItems(ctx *gin.Context) {
	items, err := s.store.ListClearanceItems(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

// GET /departments/:id/clearance-items
func (s *Server) listItemsByDepartment(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Optional: verify department exists
	_, err = s.store.GetDepartment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("department not found"))
		return
	}

	items, err := s.store.ListItemsByDepartment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, items)
}

type updateClearanceItemRequest struct {
	Code               string `json:"code" binding:"required"`
	Title              string `json:"title" binding:"required"`
	Description        string `json:"description"`
	DepartmentID       int64  `json:"department_id" binding:"required"`
	ApproverStaffID    int64  `json:"approver_staff_id" binding:"required"`
	RequiresAttachment bool   `json:"requires_attachment"`
	Sequence           int32  `json:"sequence" binding:"required,min=1"`
}

// PUT /clearance-items/:id
func (s *Server) updateClearanceItem(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateClearanceItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate references
	if _, err := s.store.GetDepartment(ctx, req.DepartmentID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid department ID"))
		return
	}

	if _, err := s.store.GetStaffUser(ctx, req.ApproverStaffID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid approver staff ID"))
		return
	}

	arg := db.UpdateClearanceItemParams{
		Code:               req.Code,
		Title:              req.Title,
		Description:        req.Description,
		DepartmentID:       req.DepartmentID,
		ApproverStaffID:    req.ApproverStaffID,
		RequiresAttachment: req.RequiresAttachment,
		Sequence:           req.Sequence,
		ID:                 id,
	}

	item, err := s.store.UpdateClearanceItem(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

// DELETE /clearance-items/:id
func (s *Server) deleteClearanceItem(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = s.store.DeleteClearanceItem(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
