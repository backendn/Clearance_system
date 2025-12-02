package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateClearanceRecordRequest struct {
	StudentID       int64  `json:"student_id" binding:"required,min=1"`
	ClearanceItemID int64  `json:"clearance_item_id" binding:"required,min=1"`
	SessionID       int64  `json:"session_id" binding:"required,min=1"`
	Note            string `json:"note"`
	AttachmentURL   string `json:"attachment_url"`
}

type UpdateClearanceRecordStatusRequest struct {
	Status        string `json:"status" binding:"required"`
	Note          string `json:"note"`
	HandledBy     int64  `json:"handled_by" binding:"required"`
	AttachmentURL string `json:"attachment_url"`
}

func (server *Server) createClearanceRecord(ctx *gin.Context) {
	var req CreateClearanceRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateClearanceRecordParams{
		StudentID:       req.StudentID,
		ClearanceItemID: req.ClearanceItemID,
		SessionID:       req.SessionID,
		Status:          "pending",
		Note:            req.Note,
		HandledBy:       0,
		AttachmentUrl:   req.AttachmentURL,
	}

	record, err := server.store.CreateClearanceRecord(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, record)
}
func (server *Server) getClearanceRecord(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	record, err := server.store.GetClearanceRecord(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, record)
}
func (server *Server) listRecordsByStudent(ctx *gin.Context) {
	studentID, err := strconv.ParseInt(ctx.Param("student_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	records, err := server.store.ListRecordsByStudent(ctx, studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, records)
}
func (server *Server) listRecordsBySession(ctx *gin.Context) {
	sessionID, err := strconv.ParseInt(ctx.Param("session_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	records, err := server.store.ListRecordsBySession(ctx, sessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, records)
}
func (server *Server) updateClearanceRecordStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateClearanceRecordStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateClearanceRecordStatusParams{
		Status:        req.Status,
		Note:          req.Note,
		HandledBy:     req.HandledBy,
		HandledAt:     time.Now(),
		AttachmentUrl: req.AttachmentURL,
		ID:            id,
	}

	record, err := server.store.UpdateClearanceRecordStatus(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// First load the clearance item
	item, err := server.store.GetClearanceItem(ctx, record.ClearanceItemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if arg.Status == "approved" {
		server.sendNotification(ctx, 0, record.StudentID,
			fmt.Sprintf("Your clearance item '%s' has been approved.", item.Title))
	} else if arg.Status == "rejected" {
		server.sendNotification(ctx, 0, record.StudentID,
			fmt.Sprintf("Your clearance item '%s' has been rejected. Note: %s", item.Title, arg.Note))
	}

	// Notify staff for confirmation
	server.sendNotification(ctx, record.HandledBy, 0,
		fmt.Sprintf("You updated clearance record %d with status '%s'.", record.ID, arg.Status))

	ctx.JSON(http.StatusOK, record)
}
func (server *Server) deleteClearanceRecord(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteClearanceRecord(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"deleted": true})
}
