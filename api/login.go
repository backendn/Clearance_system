package api

import (
	"time"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/backendn/clearance_system/token"
	"github.com/backendn/clearance_system/util"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string         `json:"access_token"`
	Payload     *token.Payload `json:"payload"`
	User        db.StaffUser   `json:"user"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req loginRequest

	// Read request JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Fetch user from DB
	user, err := server.store.GetStaffUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "user not found"})
		return
	}

	// Check password
	err = util.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "invalid password"})
		return
	}

	// Create JWT token

	role, err := server.store.GetRole(ctx, user.RoleID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "role not found"})
		return
	}

	tokenString, payload, err := server.tokenMaker.CreateToken(
		user.ID,
		role.Name,
		time.Hour,
	)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot create token"})
		return
	}

	// Return response
	resp := loginResponse{
		AccessToken: tokenString,
		Payload:     payload,
		User:        user,
	}

	ctx.JSON(200, resp)
}
