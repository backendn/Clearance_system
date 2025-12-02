package api

import (
	db "github.com/backendn/clearance_system/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for the clearance system.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routes.
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()
	server.router = router

	// Register all routes
	server.setupRoutes()

	return server
}

// setupRoutes registers all endpoints.
func (server *Server) setupRoutes() {
	// Student routes
	server.router.POST("/students", server.CreateStudent)
	server.router.GET("/students/:student_number", server.GetStudentByNumber)
	server.router.PATCH("/students/:id", server.UpdateStudent)
	server.router.DELETE("/students/:id", server.DeleteStudent)
	// Department routes
	server.router.POST("/departments", server.CreateDepartment)
	server.router.GET("/departments", server.ListDepartments)
	server.router.GET("/departments/:id", server.GetDepartment)
	server.router.PATCH("/departments/:id", server.UpdateDepartment)
	server.router.DELETE("/departments/:id", server.DeleteDepartment)
	// Staff user routes
	server.router.POST("/staff_users", server.CreateStaffUser)
	server.router.GET("/staff_users/:id", server.GetStaffUser)
	server.router.GET("/staff_users", server.ListStaffUsers)
	server.router.PATCH("/staff_users/:id", server.UpdateStaffUser)
	server.router.DELETE("/staff_users/:id", server.DeleteStaffUser)
	// clearance_items routes
	server.router.POST("/clearance_items", server.createClearanceItem)
	server.router.GET("/clearance_items/:id", server.getClearanceItem)
	server.router.GET("/clearance_items", server.listClearanceItems)
	server.router.GET("/departments/:department_id/clearance-items", server.listItemsByDepartment)
	server.router.PATCH("/clearance_items/:id", server.updateClearanceItem)
	server.router.DELETE("/clearance_items/:id", server.deleteClearanceItem)
	// Clearance Requests
	server.router.POST("/students/:id/clearance_request", server.SubmitClearanceRequest)
	server.router.GET("/students/:id/clearance_requests", server.ListStudentRequests)
	server.router.GET("/clearance_requests/:id", server.GetClearanceRequest)

	// Clearance Record routes
	server.router.POST("/clearance_records", server.createClearanceRecord)
	server.router.GET("/clearance_records/:id", server.getClearanceRecord)
	server.router.GET("/students/:student_id/records", server.listRecordsByStudent)
	server.router.GET("/sessions/:session_id/records", server.listRecordsBySession)
	server.router.PATCH("/clearance_records/:id/status", server.updateClearanceRecordStatus)
	server.router.DELETE("/clearance_records/:id", server.deleteClearanceRecord)
	// Notifications
	server.router.POST("/notifications", server.CreateNotification)
	server.router.GET("/notifications/:id", server.GetNotification)
	server.router.GET("/notifications/user/:id", server.ListNotificationsForUser)
	server.router.GET("/notifications/student/:id", server.ListNotificationsForStudent)
	server.router.PATCH("/notifications/:id/read", server.MarkNotificationRead)
	server.router.DELETE("/notifications/:id", server.DeleteNotification)

	server.router.POST("/roles", server.CreateRole)
	server.router.GET("/roles/:id", server.GetRole)
	server.router.GET("/roles", server.ListRoles)
	server.router.DELETE("/roles/:id", server.DeleteRole)

	// Add other modules later...
	// server.router.POST("/departments", server.CreateDepartment)
	// ...
}

// Start runs the HTTP server on given address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func (server *Server) sendNotification(
	ctx *gin.Context,
	userID int64,
	studentID int64,
	msg string,
) {
	arg := db.CreateNotificationParams{
		RecipientUserID:    userID,
		RecipientStudentID: studentID,
		Message:            msg,
		Read:               false,
	}

	_, err := server.store.CreateNotification(ctx, arg)
	if err != nil {
		// Basic logging, won't break workflow
		// Option 1: Gin
		ctx.Error(err)

		// Option 2: Standard log
		// log.Println("notification error:", err)

		return
	}
}
