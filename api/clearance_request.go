package api

import (
	"database/sql"
	"net/http"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/gin-gonic/gin"
)

type ClearanceRequestResponse struct {
	ID        int64  `json:"id"`
	StudentID int64  `json:"student_id"`
	SessionID int64  `json:"session_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func convertClearanceRequest(r db.ClearanceRequest) ClearanceRequestResponse {
	return ClearanceRequestResponse{
		ID:        r.ID,
		StudentID: r.StudentID,
		SessionID: r.SessionID,
		Status:    r.Status,
		CreatedAt: r.CreatedAt.String(),
	}
}
func (server *Server) SubmitClearanceRequest(ctx *gin.Context) {
	studentID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid student ID"))
		return
	}

	// 1. Validate student exists
	_, err = server.store.GetStudent(ctx, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage("student not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	// 2. Get active clearance session
	session, err := server.store.GetActiveSession(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorMessage("no active clearance session"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	// 3. Ensure no duplicate request for this session
	_, err = server.store.GetStudentRequestForSession(ctx, db.GetStudentRequestForSessionParams{
		StudentID: studentID,
		SessionID: session.ID,
	})
	if err == nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("clearance request already submitted"))
		return
	}

	// 4. Create a new clearance request
	req, err := server.store.CreateClearanceRequest(ctx, db.CreateClearanceRequestParams{
		StudentID: studentID,
		SessionID: session.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	// 5. Load all clearance items
	items, err := server.store.ListClearanceItems(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage("failed to load clearance items"))
		return
	}

	// 6. Create clearance_records for each clearance item
	for _, item := range items {
		_, err := server.store.CreateClearanceRecord(ctx, db.CreateClearanceRecordParams{
			StudentID:       studentID,
			ClearanceItemID: item.ID,
			SessionID:       session.ID,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorMessage("failed to create clearance workflow"))
			return
		}
	}

	// 7. Respond
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "clearance request submitted successfully",
		"request": convertClearanceRequest(req),
	})
}
func (server *Server) ListStudentRequests(ctx *gin.Context) {
	studentID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid student ID"))
		return
	}

	// load requests
	reqs, err := server.store.ListRequestsByStudent(ctx, studentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	// convert
	resp := make([]ClearanceRequestResponse, 0)
	for _, r := range reqs {
		resp = append(resp, convertClearanceRequest(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"requests": resp,
	})
}
func (server *Server) GetClearanceRequest(ctx *gin.Context) {
	reqID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid request ID"))
		return
	}

	req, err := server.store.GetClearanceRequest(ctx, reqID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage("clearance request not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, convertClearanceRequest(req))
}
