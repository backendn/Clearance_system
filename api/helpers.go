package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIDParam(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return 0, err
	}
	return id, nil
}
// errorMessage creates a simple JSON response with a message string
func errorMessage(msg string) gin.H {
	return gin.H{"error": msg}
}

// errorResponse wraps an error object into JSON
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
