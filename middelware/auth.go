package middleware

import (
	"net/http"
	"strings"

	"github.com/backendn/clearance_system/token"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		fields := strings.Fields(auth)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization"})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("authPayload", payload)
		ctx.Next()
	}
}
