package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) CreateRole(ctx *gin.Context) {
	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err.Error()))
		return
	}

	role, err := server.store.CreateRole(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"role": role})
}

func (server *Server) GetRole(ctx *gin.Context) {
	id, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage("invalid id"))
		return
	}

	role, err := server.store.GetRole(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage("role not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"role": role})
}
func (server *Server) ListRoles(ctx *gin.Context) {
	roles, err := server.store.ListRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}
func (server *Server) DeleteRole(ctx *gin.Context) {
    id, err := getIDParam(ctx)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorMessage("invalid id"))
        return
    }

    err = server.store.DeleteRole(ctx, id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "role deleted"})
}
