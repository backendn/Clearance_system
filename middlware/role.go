package middlware

import (
	"net/http"

	"github.com/backendn/clearance_system/token"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get("payload")

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token payload"})
			ctx.Abort()
			return
		}

		userPayload := payload.(*token.Payload)

		// Check if user's role is allowed
		for _, role := range allowedRoles {
			if userPayload.Role == role {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role permissions"})
		ctx.Abort()
	}
}
