package api

import (
	db "github.com/backendn/clearance_system/db/sqlc"
	middleware "github.com/backendn/clearance_system/middelware"
	"github.com/backendn/clearance_system/token"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for the clearance system.
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates a new HTTP server and sets up routes.
func NewServer(store db.Store) *Server {
	// Create JWT maker with a secret key (32+ chars)
	secretKey := "your-very-secure-32-char-secret-key"

	maker, err := token.NewJWTMaker(secretKey)
	if err != nil {
		panic("cannot create token maker: " + err.Error())
	}

	server := &Server{
		store:      store,
		tokenMaker: maker,
	}

	router := gin.Default()
	server.router = router

	server.setupRoutes()

	return server
}

// setupRoutes registers all endpoints.
func (server *Server) setupRoutes() {

	// Public endpoints
	server.router.POST("/login", server.Login) // you should create this
	server.router.POST("/register", server.CreateStaffUser)

	// Protected routes
	auth := server.router.Group("/")
	auth.Use(middleware.AuthMiddleware(server.tokenMaker))

	// Students
	auth.POST("/students", server.CreateStudent)
	auth.GET("/students/:student_number", server.GetStudentByNumber)
	auth.PATCH("/students/:id", server.UpdateStudent)
	auth.DELETE("/students/:id", server.DeleteStudent)

	// Departments
	auth.POST("/departments", server.CreateDepartment)
	auth.GET("/departments", server.ListDepartments)
	auth.GET("/departments/:id", server.GetDepartment)
	auth.PATCH("/departments/:id", server.UpdateDepartment)
	auth.DELETE("/departments/:id", server.DeleteDepartment)

	// Staff
	auth.GET("/staff_users/:id", server.GetStaffUser)
	auth.GET("/staff_users", server.ListStaffUsers)
	auth.PATCH("/staff_users/:id", server.UpdateStaffUser)
	auth.DELETE("/staff_users/:id", server.DeleteStaffUser)

	// Clearance Items
	auth.POST("/clearance_items", server.createClearanceItem)
	auth.GET("/clearance_items/:id", server.getClearanceItem)
	auth.GET("/clearance_items", server.listClearanceItems)
	auth.GET("/departments/:department_id/clearance-items", server.listItemsByDepartment)
	auth.PATCH("/clearance_items/:id", server.updateClearanceItem)
	auth.DELETE("/clearance_items/:id", server.deleteClearanceItem)

	// Clearance Requests
	auth.POST("/students/:id/clearance_request", server.SubmitClearanceRequest)
	auth.GET("/students/:id/clearance_requests", server.ListStudentRequests)
	auth.GET("/clearance_requests/:id", server.GetClearanceRequest)

	// Records
	auth.POST("/clearance_records", server.createClearanceRecord)
	auth.GET("/clearance_records/:id", server.getClearanceRecord)
	auth.GET("/students/:student_id/records", server.listRecordsByStudent)
	auth.GET("/sessions/:session_id/records", server.listRecordsBySession)
	auth.PATCH("/clearance_records/:id/status", server.updateClearanceRecordStatus)
	auth.DELETE("/clearance_records/:id", server.deleteClearanceRecord)

	// Notifications
	auth.POST("/notifications", server.CreateNotification)
	auth.GET("/notifications/:id", server.GetNotification)
	auth.GET("/notifications/user/:id", server.ListNotificationsForUser)
	auth.GET("/notifications/student/:id", server.ListNotificationsForStudent)
	auth.PATCH("/notifications/:id/read", server.MarkNotificationRead)
	auth.DELETE("/notifications/:id", server.DeleteNotification)

	// Roles
	auth.POST("/roles", server.CreateRole)
	auth.GET("/roles/:id", server.GetRole)
	auth.GET("/roles", server.ListRoles)
	auth.DELETE("/roles/:id", server.DeleteRole)
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
