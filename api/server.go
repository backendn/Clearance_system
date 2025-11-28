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

	// Add other modules later...
	// server.router.POST("/departments", server.CreateDepartment)
	// ...
}

// Start runs the HTTP server on given address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
