package api

import (
	"database/sql"
	"net/http"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

type NotificationResponse struct {
	ID                 int64  `json:"id"`
	RecipientUserID    int64  `json:"recipient_user_id"`
	RecipientStudentID int64  `json:"recipient_student_id"`
	Message            string `json:"message"`
	Read               bool   `json:"read"`
	CreatedAt          string `json:"created_at"`
}

func convertNotification(n db.Notification) NotificationResponse {
	return NotificationResponse{
		ID:                 n.ID,
		RecipientUserID:    n.RecipientUserID,
		RecipientStudentID: n.RecipientStudentID,
		Message:            n.Message,
		Read:               n.Read,
		CreatedAt:          n.CreatedAt.String(),
	}
}

type createNotificationRequest struct {
	RecipientUserID    int64  `json:"recipient_user_id"`
	RecipientStudentID int64  `json:"recipient_student_id"`
	Message            string `json:"message" binding:"required"`
}

// ================================
// Create Notification
// ================================
func (server *Server) CreateNotification(ctx *gin.Context) {
	var req createNotificationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err.Error()))
		return
	}

	notification, err := server.store.CreateNotification(ctx, db.CreateNotificationParams{
		RecipientUserID:    req.RecipientUserID,
		RecipientStudentID: req.RecipientStudentID,
		Message:            req.Message,
		Read:               false,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, convertNotification(notification))
}

// ================================
// Get Single Notification
// ================================
func (server *Server) GetNotification(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid notification ID"))
		return
	}

	n, err := server.store.GetNotification(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage("notification not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, convertNotification(n))
}

// ================================
// List Notifications For Staff User
// ================================
func (server *Server) ListNotificationsForUser(ctx *gin.Context) {
	userID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid user ID"))
		return
	}

	list, err := server.store.ListNotificationsForUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	resp := make([]NotificationResponse, 0)
	for _, n := range list {
		resp = append(resp, convertNotification(n))
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": resp})
}

// ================================
// List Notifications For Student
// ================================
func (server *Server) ListNotificationsForStudent(ctx *gin.Context) {
	studentID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid student ID"))
		return
	}

	list, err := server.store.ListNotificationsForStudent(ctx, studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	resp := make([]NotificationResponse, 0)
	for _, n := range list {
		resp = append(resp, convertNotification(n))
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": resp})
}

// ================================
// Mark Notification as Read
// ================================
func (server *Server) MarkNotificationRead(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid notification ID"))
		return
	}

	n, err := server.store.MarkNotificationRead(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage("notification not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, convertNotification(n))
}

// ================================
// Delete Notification
// ================================
func (server *Server) DeleteNotification(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid notification ID"))
		return
	}

	err = server.store.DeleteNotification(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
}
