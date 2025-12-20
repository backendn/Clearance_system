package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
func getPagination(ctx *gin.Context) (limit, offset int) {
	limitQuery := ctx.DefaultQuery("limit", "10")
	offsetQuery := ctx.DefaultQuery("offset", "0")

	fmt.Sscanf(limitQuery, "%d", &limit)
	fmt.Sscanf(offsetQuery, "%d", &offset)

	if limit < 1 {
		limit = 10
	}
	return
}
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}

// CheckPassword compares a hashed password with the plain password.
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NullableString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func NullIntToInt(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}
func ToNullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: v != 0}
}
