package middlware

import (
	"net/http"

	"github.com/backendn/clearance_system/token"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet("authorization_payload").(*token.Payload)

		if payload.Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			return
		}

		ctx.Next()
	}
}
