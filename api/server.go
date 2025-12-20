package api

import (
	db "github.com/backendn/clearance_system/db/sqlc"
	middleware "github.com/backendn/clearance_system/middelware"
	"github.com/backendn/clearance_system/token"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates a new HTTP server and configures routes
func NewServer(store db.Store) *Server {
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

func (server *Server) setupRoutes() {

	// --------------------
	// PUBLIC ROUTES
	// --------------------
	server.router.POST("/login", server.Login)
	server.router.POST("/register", server.CreateStaffUser) // only for now

	// --------------------
	// AUTHENTICATED ROUTES
	// --------------------
	auth := server.router.Group("/")
	auth.Use(middleware.AuthMiddleware(server.tokenMaker))

	// --------------------
	// ADMIN ONLY
	// --------------------
	admin := auth.Group("/")
	admin.Use(middleware.RoleMiddleware("admin"))

	admin.POST("/departments", server.CreateDepartment)
	admin.DELETE("/departments/:id", server.DeleteDepartment)

	admin.POST("/roles", server.CreateRole)
	admin.DELETE("/roles/:id", server.DeleteRole)

	admin.POST("/clearance_items", server.createClearanceItem)
	admin.DELETE("/staff_users/:id", server.DeleteStaffUser)

	// --------------------
	// STAFF ONLY
	// --------------------
	staff := auth.Group("/")
	staff.Use(middleware.RoleMiddleware("staff", "admin"))

	staff.PATCH("/clearance_records/:id/status", server.updateClearanceRecordStatus)
	staff.GET("/sessions/:session_id/records", server.listRecordsBySession)

	// --------------------
	// STUDENT ONLY
	// --------------------
	student := auth.Group("/")
	student.Use(middleware.RoleMiddleware("student"))

	student.POST("/students/:id/clearance_request", server.SubmitClearanceRequest)
	student.GET("/students/:id/clearance_requests", server.ListStudentRequests)

	// --------------------
	// GENERAL AUTH ROUTES (everyone with login)
	// --------------------

	// Students
	auth.POST("/students", server.CreateStudent)
	auth.GET("/students/number/:student_number", server.GetStudentByNumber)
	auth.PATCH("/students/:id", server.UpdateStudent)
	auth.DELETE("/students/:id", server.DeleteStudent)

	// Departments
	auth.GET("/departments", server.ListDepartments)
	auth.GET("/departments/:id", server.GetDepartment)
	auth.PATCH("/departments/:id", server.UpdateDepartment)

	// Staff
	auth.GET("/staff_users/:id", server.GetStaffUser)
	auth.GET("/staff_users", server.ListStaffUsers)
	auth.PATCH("/staff_users/:id", server.UpdateStaffUser)

	// Clearance Items
	auth.GET("/clearance_items", server.listClearanceItems)
	auth.GET("/clearance_items/:id", server.getClearanceItem)
	auth.GET("/departments/department/:department_id/clearance-items", server.listItemsByDepartment)
	auth.PATCH("/clearance_items/:id", server.updateClearanceItem)
	auth.DELETE("/clearance_items/:id", server.deleteClearanceItem)

	// Clearance Requests
	auth.GET("/clearance_requests/:id", server.GetClearanceRequest)

	// Records
	auth.POST("/clearance_records", server.createClearanceRecord)
	auth.GET("/clearance_records/:id", server.getClearanceRecord)
	auth.GET("/students/student/:student_id/records", server.listRecordsByStudent)
	auth.DELETE("/clearance_records/:id", server.deleteClearanceRecord)

	// Notifications
	auth.POST("/notifications", server.CreateNotification)
	auth.GET("/notifications/:id", server.GetNotification)
	auth.GET("/notifications/user/:id", server.ListNotificationsForUser)
	auth.GET("/notifications/student/:id", server.ListNotificationsForStudent)
	auth.PATCH("/notifications/:id/read", server.MarkNotificationRead)
	auth.DELETE("/notifications/:id", server.DeleteNotification)

}

// Start server
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
		RecipientUserID:    ToNullInt64(userID),
		RecipientStudentID: ToNullInt64(studentID),
		Message:            msg,
		Read:               false,
	}

	_, err := server.store.CreateNotification(ctx, arg)
	if err != nil {
		// Basic logging, won't break workflow
		ctx.Error(err) // logs to Gin

		// OR standard logging:
		// log.Println("notification error:", err)

		return
	}
}
