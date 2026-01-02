package api

import (
	"net/http"
	"time"

	db "github.com/backendn/clearance_system/db/sqlc"
	"github.com/backendn/clearance_system/util"
	"github.com/gin-gonic/gin"
)

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func (server *Server) LoginAdmin(ctx *gin.Context) {
	var req LoginAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	admin, err := server.store.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorMessage("invalid credentials"))
		return
	}

	if err := util.CheckPassword(req.Password, admin.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorMessage("invalid credentials"))
		return
	}

	tokenStr, _, err := server.tokenMaker.CreateToken(
		admin.ID,
		"admin",
		time.Hour*24,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage("failed to create token"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": tokenStr,
		"admin": AdminResponse{
			ID:        admin.ID,
			Username:  admin.Username,
			FullName:  admin.FullName,
			Email:     admin.Email,
			Role:      admin.Role,
			IsActive:  admin.IsActive,
			CreatedAt: admin.CreatedAt.Format(time.RFC3339),
		},
	})
}
func (server *Server) CreateAdmin(ctx *gin.Context) {
	var req CreateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage("failed to hash password"))
		return
	}

	admin, err := server.store.CreateAdmin(ctx, db.CreateAdminParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		Role:           "admin",
		IsActive:       true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		FullName:  admin.FullName,
		Email:     admin.Email,
		Role:      admin.Role,
		IsActive:  admin.IsActive,
		CreatedAt: admin.CreatedAt.Format(time.RFC3339),
	})
}
func (server *Server) ListAdmins(ctx *gin.Context) {
	admins, err := server.store.ListAdmins(ctx, db.ListAdminsParams{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	var res []AdminResponse
	for _, a := range admins {
		res = append(res, AdminResponse{
			ID:        a.ID,
			Username:  a.Username,
			FullName:  a.FullName,
			Email:     a.Email,
			Role:      a.Role,
			IsActive:  a.IsActive,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		})
	}

	ctx.JSON(http.StatusOK, res)
}
